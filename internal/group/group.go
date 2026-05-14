package group

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"gochat/internal/infra/db"
	"gochat/internal/model"
	"gochat/internal/protocol"
	"strings"
)

var (
	ErrGroupDismissed    = errors.New("group is dismissed")
	ErrGroupNotOpen      = errors.New("group is not open to join")
	ErrLevelInsufficient = errors.New("user level insufficient")
	ErrGroupFull         = errors.New("group is full")
	ErrAlreadyMember     = errors.New("already a member")
	ErrSiteMismatch      = errors.New("site mismatch")
)

const (
	RoleMember = int8(0)
	RoleMod    = int8(1)
	RoleOwner  = int8(2)
)

func genCode(siteBid string) string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	prefix := strings.ToUpper(siteBid)
	if len(prefix) > 4 {
		prefix = prefix[:4]
	}
	return prefix + hex.EncodeToString(b)
}

const getGroupsOfUserSQL = `
SELECT cg.id, cg.title, cg.code, cg.open_join, cg.join_user_level,
       COALESCE(lr.content, '') AS last_msg,
       COALESCE(EXTRACT(EPOCH FROM lr.create_time)::bigint, 0) AS last_msg_time
FROM chat_group cg
JOIN group_user gu ON gu.group_id = cg.id
    AND gu.user_id = $1
    AND gu.deleted = false
LEFT JOIN LATERAL (
    SELECT content, create_time
    FROM chat_record
    WHERE group_id = cg.id AND deleted = false
    ORDER BY id DESC LIMIT 1
) lr ON true
WHERE cg.site_bid = $2
  AND cg.is_dismiss = false`

func Create(req *protocol.CreateGroupReq, owner *model.User) (int64, string, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return 0, "", err
	}

	code := req.Code
	if code == "" {
		code = genCode(owner.SiteBid)
	}

	tx, err := d.DB.Beginx()
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback()

	var groupID int64
	err = d.Builder.
		Insert("chat_group").
		Columns("site_bid", "title", "code", "open_join", "user_limit", "bulletin",
			"owner_user_id", "owner_ext_member_id", "owner_ext_username").
		Values(owner.SiteBid, req.Title, code, req.OpenJoin, req.UserLimit, req.Bulletin,
			owner.ID, owner.ExtMemberID, owner.ExtUsername).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().
		Scan(&groupID)
	if err != nil {
		return 0, "", err
	}

	_, err = d.Builder.
		Insert("group_user").
		Columns("user_id", "group_id", "role_type").
		Values(owner.ID, groupID, RoleOwner).
		RunWith(tx).
		Exec()
	if err != nil {
		return 0, "", err
	}

	return groupID, code, tx.Commit()
}

func GetMembership(userID int64, groupID int64) (*model.GroupUser, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	querySql, args, _ := d.Builder.
		Select("id", "user_id", "group_id", "role_type", "is_ban", "deleted", "last_visible_chat_id", "last_read_chat_id").
		From("group_user").
		Where(sq.Eq{"user_id": userID, "group_id": groupID, "deleted": false}).
		Limit(1).ToSql()

	var gu model.GroupUser
	if err := d.DB.Get(&gu, querySql, args...); err != nil {
		return nil, err
	}

	return &gu, nil
}

func GetMemberCount(groupID int64) (int, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return 0, err
	}

	var count int
	err = d.Builder.
		Select("COUNT(*)").
		From("group_user").
		Where(sq.Eq{"group_id": groupID, "deleted": false}).
		RunWith(d.DB).QueryRow().Scan(&count)
	return count, err
}

var groupColumns = []string{
	"id", "site_bid", "title", "code", "is_dismiss", "open_join",
	"user_limit", "join_user_level", "speak_user_level", "owner_user_id", "visible",
}

func FindByID(groupID int64) (*model.ChatGroup, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	querySql, args, _ := d.Builder.
		Select(groupColumns...).
		From("chat_group").
		Where(sq.Eq{"id": groupID}).
		Limit(1).ToSql()

	var g model.ChatGroup
	if err := d.DB.Get(&g, querySql, args...); err != nil {
		return nil, err
	}

	return &g, nil
}

func Join(u *model.User, g *model.ChatGroup) error {
	if g.IsDismiss {
		return ErrGroupDismissed
	}
	if u.SiteBid != g.SiteBid {
		return ErrSiteMismatch
	}
	if !g.OpenJoin {
		return ErrGroupNotOpen
	}
	if g.JoinUserLevel > 0 && u.UserLevel < g.JoinUserLevel {
		return ErrLevelInsufficient
	}

	d, err := db.GetDBConn()
	if err != nil {
		return err
	}

	if g.UserLimit > 0 {
		count, err := GetMemberCount(g.ID)
		if err != nil {
			return err
		}
		if count >= g.UserLimit {
			return ErrGroupFull
		}
	}

	var existing struct {
		ID      int64 `db:"id"`
		Deleted bool  `db:"deleted"`
	}
	querySql, args, _ := d.Builder.
		Select("id", "deleted").
		From("group_user").
		Where(sq.Eq{"user_id": u.ID, "group_id": g.ID}).
		Limit(1).ToSql()

	err = d.DB.Get(&existing, querySql, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil {
		if !existing.Deleted {
			return ErrAlreadyMember
		}
		_, err = d.Builder.Update("group_user").
			Set("deleted", false).
			Where(sq.Eq{"id": existing.ID}).
			RunWith(d.DB).Exec()
		return err
	}

	_, err = d.Builder.
		Insert("group_user").
		Columns("user_id", "group_id", "role_type").
		Values(u.ID, g.ID, RoleMember).
		RunWith(d.DB).Exec()
	return err
}

func GetGroupsOfUser(userID int64, siteBid string) ([]protocol.DisplayUserGroup, error) {
	d, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}

	rows := make([]protocol.DisplayUserGroup, 0)
	if err := d.DB.Select(&rows, getGroupsOfUserSQL, userID, siteBid); err != nil {
		return nil, err
	}

	return rows, nil
}
