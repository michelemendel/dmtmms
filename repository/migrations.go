package repo

import (
	"fmt"
	"log/slog"

	"time"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

func (r *Repo) runStatements(sqlStmts map[string]string) {
	for name, sqlStmt := range sqlStmts {
		_, err := r.DB.Exec(sqlStmt)
		if err != nil {
			slog.Error(fmt.Sprintf("Error in stmt %q:%s\n%s\n", err, name, sqlStmt))
		}
	}
}

// --------------------------------------------------------------------------------
// DDL
func (r *Repo) DropTables() {
	var sqlStmts = make(map[string]string)

	// Tables
	sqlStmts["drop_users"] = `DROP TABLE IF EXISTS users;`
	sqlStmts["drop_members"] = `DROP TABLE IF EXISTS members;`
	sqlStmts["drop_families"] = `DROP TABLE IF EXISTS families;`
	sqlStmts["drop_groups"] = `DROP TABLE IF EXISTS groups;`
	sqlStmts["drop_members_groups"] = `DROP TABLE IF EXISTS members_groups;`

	// Triggers
	sqlStmts["drop_trigger_inc_member_id"] = `DROP TRIGGER IF EXISTS trigger_inc_member_id;`
	sqlStmts["drop_trigger_u_members_updated_at"] = `DROP TRIGGER IF EXISTS trigger_u_members_updated_at;`
	sqlStmts["drop_trigger_u_families_updated_at"] = `DROP TRIGGER IF EXISTS trigger_u_families_updated_at;`
	sqlStmts["drop_trigger_u_users_updated_at"] = `DROP TRIGGER IF EXISTS trigger_u_users_updated_at;`
	sqlStmts["drop_trigger_u_members_groups_updated_at"] = `DROP TRIGGER IF EXISTS trigger_u_members_groups_updated_at;`

	// Indexes
	sqlStmts["drop_index_members_uuid"] = `DROP INDEX IF EXISTS idx_members_uuid;`
	sqlStmts["drop_index_families_uuid"] = `DROP INDEX IF EXISTS idx_families_uuid;`
	sqlStmts["drop_index_groups_uuid"] = `DROP INDEX IF EXISTS idx_groups_uuid;`
	sqlStmts["drop_index_members_groups"] = `DROP INDEX IF EXISTS idx_members_groups;`

	r.runStatements(sqlStmts)
}

func (r *Repo) CreateTables() {
	var sqlStmts = make(map[string]string)

	sqlStmts["create_users"] = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	); `

	// id - this is an id shared by dmt and tripletex
	// https://mail.google.com/mail/u/0/#search/from%3Aadm%40dmt.oslo.no/FMfcgzGwJcbCTnSbdzxfhHXlFLZCXDGZ?projector=1&messagePartId=0.1
	// state = [active, archived, tobedeleted]
	sqlStmts["create_members"] = `
	CREATE TABLE IF NOT EXISTS members (
		uuid TEXT PRIMARY KEY,
		id INTEGER UNIQUE,
		name TEXT NOT NULL,
		dob REAL,
		personnummer TEXT,
		email TEXT,
		mobile TEXT,
		address1 TEXT,
		address2 TEXT,
		postnummer TEXT,
		poststed TEXT,
		synagogue_seat TEXT,
		membership_fee_tier TEXT,
		registered_date REAL,
		deregistered_date REAL,
		receive_email BOOLEAN,
		receive_mail BOOLEAN,
		receive_hatikvah BOOLEAN,
		status TEXT,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		family_uuid TEXT,
		family_name TEXT,
		FOREIGN KEY(family_uuid) REFERENCES families(uuid)
	); `

	sqlStmts["create_families"] = `
	CREATE TABLE IF NOT EXISTS families (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	sqlStmts["create_groups"] = `
	CREATE TABLE IF NOT EXISTS groups (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	// many-to-many between members and groups
	sqlStmts["create_members_groups"] = `
	CREATE TABLE IF NOT EXISTS members_groups (
		member_uuid TEXT NOT NULL,
		group_uuid TEXT NOT NULL,
		role TEXT,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		primary key (member_uuid, group_uuid)
		FOREIGN KEY(member_uuid) REFERENCES members(uuid),
		FOREIGN KEY(group_uuid) REFERENCES groups(uuid)
		); `
	// FOREIGN KEY(member_uuid) REFERENCES members(uuid) ON UPDATE CASCADE ON DELETE CASCADE,
	// FOREIGN KEY(group_uuid) REFERENCES groups(uuid)

	r.runStatements(sqlStmts)
}

// SELECT name FROM sqlite_master WHERE type = 'trigger';
func (r *Repo) CreateTriggers() {
	var sqlStmts = make(map[string]string)

	// CREATE TRIGGER inc_member_id AFTER INSERT ON names WHEN new.id IS NULL BEGIN UPDATE names SET id=(SELECT coalesce(max(id),0)+1 FROM names) WHERE uuid=new.uuid; END;

	sqlStmts["trigger_inc_member_id"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_inc_member_id AFTER INSERT ON members
	WHEN new.id IS NULL BEGIN 
		UPDATE members SET id=(SELECT coalesce(max(id),0)+1 FROM members) 
			WHERE uuid=new.uuid; 
	END;`

	sqlStmts["trigger_u_members_updated_at"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_u_member_updated_at AFTER UPDATE ON members
	BEGIN 
		UPDATE members SET updated_at=CURRENT_TIMESTAMP WHERE uuid=new.uuid; 
	END;`

	sqlStmts["trigger_u_families_updated_at"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_u_families_updated_at AFTER UPDATE ON families
	BEGIN 
		UPDATE families SET updated_at=CURRENT_TIMESTAMP WHERE uuid=new.uuid; 
	END;`

	sqlStmts["trigger_u_groups_updated_at"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_u_groups_updated_at AFTER UPDATE ON groups
	BEGIN 
		UPDATE groups SET updated_at=CURRENT_TIMESTAMP WHERE uuid=new.uuid; 
	END;`

	sqlStmts["trigger_u_members_groups_updated_at"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_u_members_groups_updated_at AFTER UPDATE ON members_groups
	BEGIN 
		UPDATE members_groups SET updated_at=CURRENT_TIMESTAMP WHERE uuid=new.uuid; 
	END;`

	sqlStmts["trigger_u_users_updated_at"] = `
	CREATE TRIGGER IF NOT EXISTS trigger_u_users_updated_at AFTER UPDATE ON users
	BEGIN 
		UPDATE users SET updated_at=CURRENT_TIMESTAMP WHERE name=new.name; 
	END;`

	r.runStatements(sqlStmts)
}

// SELECT name FROM sqlite_master WHERE type = 'index';
func (r *Repo) CreateIndexes() {
	var sqlStmts = make(map[string]string)

	sqlStmts["create_index_members_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_members_uuid ON members(uuid);`

	sqlStmts["create_index_families_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_families_uuid ON families(uuid);`

	sqlStmts["create_index_groups_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_groups_uuid ON groups(uuid);`

	sqlStmts["create_index_members_groups"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_members_groups ON members_groups(member_uuid, group_uuid);`

	r.runStatements(sqlStmts)
}

//--------------------------------------------------------------------------------
// DML

// Users
func (r *Repo) InsertUsers() {
	users := []entity.User{
		{Name: "root", HashedPassword: "shoreshdeep", Role: "root"},
		{Name: "admin", HashedPassword: "Dimethyltryptamine", Role: "admin"},
		{Name: "eve", HashedPassword: "eve123", Role: "edit"},
		{Name: "ron", HashedPassword: "ron123", Role: "read"},
	}

	for _, user := range users {
		hpw, _ := util.HashPassword(user.HashedPassword)
		user.HashedPassword = hpw
		err := r.CreateUser(user)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Families
func (r *Repo) InsertFamilies() {
	families := []entity.Family{
		{UUID: "0", Name: "none"},
		// {UUID: "101", Name: "Cohen"},
		// {UUID: "102", Name: "Levi"},
		// {UUID: "103", Name: "Israel"},
		// {UUID: "104", Name: "Hoffman"},
	}

	for _, family := range families {
		err := r.CreateFamily(family)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Groups
func (r *Repo) InsertGroups() {
	groups := []entity.Group{
		{UUID: "0", Name: "none"},
		// {UUID: "1001", Name: "styret"},
		// {UUID: "1002", Name: "chevre"},
		// {UUID: "1003", Name: "kiddush"},
		// {UUID: "1004", Name: "ligning"},
		// {UUID: "1005", Name: "lærer"},
		// {UUID: "1006", Name: "barnehage"},
	}
	for _, group := range groups {
		err := r.CreateGroup(group)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Members
type member struct {
	name         string
	receiveEmail bool
	// archived     bool
	familyUUID string
	familyName string
	groupUUIDs []string
}

func getMembers(nofMembers int) map[string]member {
	members := map[string]member{
		"11": {name: "abe", receiveEmail: true, familyUUID: "101", familyName: "fam1", groupUUIDs: []string{"1001", "1002"}},
		"12": {name: "bob", receiveEmail: true, familyUUID: "101", familyName: "fam1", groupUUIDs: []string{"1002", "1003"}},
		"13": {name: "carl", receiveEmail: false, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1001", "1002"}},
		// "14": {name: "dave", receiveEmail: true, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1002", "1003"}},
		// "15": {name: "eve", receiveEmail: false, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1002", "1004"}},
	}

	// for i := 16; i < nofMembers; i++ {
	// 	members[fmt.Sprintf("%d", i)] = member{name: fmt.Sprintf("m%d", i), receiveEmail: false, familyUUID: "", familyName: "", groupUUIDs: []string{}}
	// }

	return members
}

func (r *Repo) InsertMembers() {
	// First member is a special case,since we are using triggers on id, and we need a start value.
	r.InitMember()

	// Create members
	dob := util.String2Date("1980-02-01")
	personnummer := "12345"
	phonePrefix := "12377"
	var status entity.MemberStatus = entity.MemberStatusActive
	for i, m := range getMembers(50) {
		memberUUID := util.GenerateUUID()
		email := m.name + "@test.com"
		mobile := phonePrefix + i

		member := entity.NewMember(
			memberUUID, 0, m.name, dob, 0, personnummer, email, mobile,
			entity.Address{}, "", "", util.String2Date("2020-01-01"), time.Time{},
			m.receiveEmail, false, false, status,
			m.familyUUID, m.familyName,
			time.Time{}, time.Time{},
		)
		err := r.CreateMember(member, m.groupUUIDs)
		if err != nil {
			slog.Error(err.Error())
		}
		dob = dob.AddDate(0, 1, 1)
	}
}

func (r *Repo) ShowUsers() error {
	var name string
	var password string
	var role string
	var createdAt string
	// -- AND dob < julianday('1965-07-22')
	// rows, err := r.DB.Query(`SELECT name,date(dob),datetime(created_at,'LOCALTIME') FROM users;`)
	rows, err := r.DB.Query(`SELECT name,password,role,datetime(created_at,'LOCALTIME') FROM users;`)
	if err != nil {
		slog.Error("error getting users", "error", err.Error())
		return err
	}
	defer rows.Close()
	frmt := "%s\t%s\t%s\t%s\n"
	fmt.Printf(frmt, "name", "password", "role", "createdAt")
	for rows.Next() {
		err := rows.Scan(&name, &password, &role, &createdAt)
		if err != nil {
			slog.Error(err.Error())
		}
		fmt.Printf(frmt, name, password, role, createdAt)
	}

	return nil
}
