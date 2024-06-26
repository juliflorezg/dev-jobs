package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/juliflorezg/dev-jobs/internal/validator"
)

var DefaultCompanyIcon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 512" width="50" height="50"><path d="M384 96V320H64L64 96H384zM64 32C28.7 32 0 60.7 0 96V320c0 35.3 28.7 64 64 64H181.3l-10.7 32H96c-17.7 0-32 14.3-32 32s14.3 32 32 32H352c17.7 0 32-14.3 32-32s-14.3-32-32-32H277.3l-10.7-32H384c35.3 0 64-28.7 64-64V96c0-35.3-28.7-64-64-64H64zm464 0c-26.5 0-48 21.5-48 48V432c0 26.5 21.5 48 48 48h64c26.5 0 48-21.5 48-48V80c0-26.5-21.5-48-48-48H528zm16 64h32c8.8 0 16 7.2 16 16s-7.2 16-16 16H544c-8.8 0-16-7.2-16-16s7.2-16 16-16zm-16 80c0-8.8 7.2-16 16-16h32c8.8 0 16 7.2 16 16s-7.2 16-16 16H544c-8.8 0-16-7.2-16-16zm32 160a32 32 0 1 1 0 64 32 32 0 1 1 0-64z" fill="#8a8eb4"/></svg>`

var DefaultCompanyBgColor = `hsl(234, 100%, 94%)`

type JobPostModelInterface interface {
	Latest() ([]JobPost, error)
	FilterPosts(position, location, contract string) ([]JobPost, error)
	Get(id int) (JobPost, error)
	InsertCompany(name, logoSVG, logoBgColor, website string) error
	InsertJobPost(companyUserID int, jobPostData CreateJopPostFields) error
	DeleteJobPost(jobPostID int) error
	EditJobPost(jobPostData EditJopPostFields) error
	InsertJobApplication(jobPostID, userID int, cvName, cvPath, coverLetter string) error
}

type JobPost struct {
	ID           int
	Position     string
	Description  string
	Contract     string
	Location     string
	PostedAt     time.Time
	Company      Company
	Requirements Requirements
	Role         Role
}

type Company struct {
	CompanyID   int
	Name        string
	LogoSVG     string
	LogoBgColor string
	Website     string
}

type Requirements struct {
	ReqID                   int
	RequirementsDescription string
	RequirementsList        []string
}

type Role struct {
	RoleID          int
	RoleDescription string
	RoleList        []string
}

type JobApplication struct {
	ID          int
	jobPostID   int
	userID      int
	cvName      string
	cvPath      string
	coverLetter string
	appliedAt   time.Time
}

type JobPostModel struct {
	DB *sql.DB
}

type CreateJopPostFields struct {
	Position     string
	Description  string
	Contract     string
	Location     string
	Requirements struct {
		Content string
		Items   []string
	}
	Role struct {
		Content string
		Items   []string
	}
	validator.Validator
}

type EditJopPostFields struct {
	ID           int
	Position     string
	Description  string
	Contract     string
	Location     string
	CompanyID    int
	Requirements struct {
		ReqID   int
		Content string
		Items   []string
	}
	Role struct {
		RoleID  int
		Content string
		Items   []string
	}
	validator.Validator
}

func (jp *JobPostModel) Latest() ([]JobPost, error) {
	// return nil, nil

	// stmt := `SELECT jp.job_post_id, jp.position, jp.description, jp.contract, jp.location, jp.posted_at, cp.name, cp.logo_svg, cp.logo_bg_color, cp.website, rq.requirements_description, rq.requirements_list, rl.role_description, rl.role_list
	// FROM jobposts AS jp
	// INNER JOIN companies AS cp ON jp.company_id = cp.company_id
	// INNER JOIN requirements AS rq ON jp.requirements_id = rq.req_id
	// INNER JOIN roles AS rl ON jp.role_id = rl.role_id
	// ORDER BY jp.posted_at LIMIT 10`

	stmt := `SELECT jp.job_post_id, jp.position, jp.description, jp.contract, jp.location, jp.posted_at, cp.name, cp.logo_svg, cp.logo_bg_color
					FROM jobposts AS jp 
					INNER JOIN companies AS cp ON jp.company_id = cp.company_id
					ORDER BY jp.posted_at DESC LIMIT 10`

	rows, err := jp.DB.Query(stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobposts []JobPost

	for rows.Next() {
		var jp JobPost
		// var rqList any
		// var roleList any

		// err := rows.Scan(&jp.ID, &jp.Position, &jp.Description, &jp.Contract, &jp.Location, &jp.PostedAt, &jp.Company.Name, &jp.Company.LogoSVG, &jp.Company.LogoBgColor, &jp.Company.Website, &jp.Requirements.RequirementsDescription, &rqList, &jp.Role.RoleDescription, &roleList)
		err := rows.Scan(&jp.ID, &jp.Position, &jp.Description, &jp.Contract, &jp.Location, &jp.PostedAt, &jp.Company.Name, &jp.Company.LogoSVG, &jp.Company.LogoBgColor)
		if err != nil {
			return nil, err
		}

		// err = json.Unmarshal(rqList.([]byte), &jp.Requirements.RequirementsList)
		// if err != nil {
		// 	return nil, err
		// }
		// err = json.Unmarshal(roleList.([]byte), &jp.Role.RoleList)
		// if err != nil {
		// 	return nil, err
		// }

		jobposts = append(jobposts, jp)
	}

	return jobposts, nil
}

func (jp *JobPostModel) FilterPosts(position, location, contract string) ([]JobPost, error) {
	stmt := `SELECT jp.job_post_id, jp.position, jp.description, jp.contract, jp.location, jp.posted_at, cp.name, cp.logo_svg, cp.logo_bg_color
			FROM jobposts AS jp 
			INNER JOIN companies AS cp ON jp.company_id = cp.company_id WHERE 1=1`

	params := make([]any, 0, 3)

	if position != "" {
		stmt += ` AND jp.position LIKE ?`
		params = append(params, "%"+position+"%")
	}
	if location != "" && location != "Filter by location" {
		stmt += ` AND jp.location LIKE ?`
		params = append(params, "%"+location+"%")
	}
	if contract != "" {
		stmt += ` AND jp.contract LIKE ?`
		params = append(params, "%"+contract+"%")
	}
	fmt.Println("params ==>", params)

	stmt += ` ORDER BY jp.posted_at DESC LIMIT 10`
	fmt.Println("sql statement ==>", stmt)

	rows, err := jp.DB.Query(stmt, params...)

	if err != nil {
		// return nil, err
		if errors.Is(err, sql.ErrNoRows) {
			return []JobPost{}, ErrNoRecord
		} else {
			return []JobPost{}, err
		}
	}
	defer rows.Close()

	var jobposts []JobPost

	for rows.Next() {
		var jp JobPost

		err := rows.Scan(&jp.ID, &jp.Position, &jp.Description, &jp.Contract, &jp.Location, &jp.PostedAt, &jp.Company.Name, &jp.Company.LogoSVG, &jp.Company.LogoBgColor)
		if err != nil {
			return nil, err
		}

		jobposts = append(jobposts, jp)
	}

	return jobposts, nil
}

func (jp *JobPostModel) Get(id int) (JobPost, error) {

	stmt := `SELECT jp.job_post_id, jp.position, jp.description, jp.contract, jp.location, jp.posted_at, cp.company_id, cp.name, cp.logo_svg, cp.logo_bg_color, cp.website, rq.req_id, rq.requirements_description, rq.requirements_list, rl.role_id, rl.role_description, rl.role_list
		FROM jobposts AS jp
		INNER JOIN companies AS cp ON jp.company_id = cp.company_id
		INNER JOIN requirements AS rq ON jp.requirements_id = rq.req_id
		INNER JOIN roles AS rl ON jp.role_id = rl.role_id
		WHERE jp.job_post_id = ?`

	row := jp.DB.QueryRow(stmt, id)
	var jobPost JobPost

	var rqList any
	var roleList any

	err := row.Scan(&jobPost.ID, &jobPost.Position, &jobPost.Description, &jobPost.Contract, &jobPost.Location, &jobPost.PostedAt, &jobPost.Company.CompanyID, &jobPost.Company.Name, &jobPost.Company.LogoSVG, &jobPost.Company.LogoBgColor, &jobPost.Company.Website, &jobPost.Requirements.ReqID, &jobPost.Requirements.RequirementsDescription, &rqList, &jobPost.Role.RoleID, &jobPost.Role.RoleDescription, &roleList)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return JobPost{}, ErrNoRecord
		} else {
			return JobPost{}, err
		}
	}

	err = json.Unmarshal(rqList.([]byte), &jobPost.Requirements.RequirementsList)
	if err != nil {
		return JobPost{}, err
	}
	err = json.Unmarshal(roleList.([]byte), &jobPost.Role.RoleList)
	if err != nil {
		return JobPost{}, err
	}

	return jobPost, nil
}

func (jp *JobPostModel) InsertCompany(name, logoSVG, logoBgColor, website string) error {
	return nil
}

func (jp *JobPostModel) InsertJobPost(companyUserID int, jobPostData CreateJopPostFields) error {
	stmt1 := `SELECT company_id FROM users_employers WHERE user_id = ?`
	var companyID int

	row := jp.DB.QueryRow(stmt1, companyUserID)

	err := row.Scan(&companyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	stmt2 := `SELECT EXISTS(SELECT true FROM companies WHERE company_id = ?)`
	var exists bool

	err = jp.DB.QueryRow(stmt2, companyID).Scan(&exists)
	if err != nil {
		return ErrNoCompany
	}

	stmt3 := `INSERT INTO requirements(requirements_description,
	requirements_list) VALUES (?, ?)`

	reqListJSON, err := json.Marshal(jobPostData.Requirements.Items)
	if err != nil {
		return ErrCouldNotConvertToJSON
	} else {
		fmt.Println()
		fmt.Printf("JSON value of req list:: %s \n", reqListJSON)
		fmt.Println()
	}

	insertedReq, err := jp.DB.Exec(stmt3, jobPostData.Requirements.Content, reqListJSON)
	if err != nil {
		return err
	}

	lastReqInsertedID, err := insertedReq.LastInsertId()
	if err != nil {
		return err
	}

	/////////////////////////////////////////////////////////////
	stmt4 := `INSERT INTO roles(role_description,
	role_list) VALUES (?, ?)`

	roleListJSON, err := json.Marshal(jobPostData.Role.Items)
	if err != nil {
		return ErrCouldNotConvertToJSON
	} else {
		fmt.Println()
		fmt.Printf("JSON value of role list:: %s \n", roleListJSON)
		fmt.Println()
	}

	insertedRole, err := jp.DB.Exec(stmt4, jobPostData.Role.Content, roleListJSON)
	if err != nil {
		return err
	}

	lastRoleInsertedID, err := insertedRole.LastInsertId()
	if err != nil {
		return err
	}

	stmt5 := `INSERT INTO jobposts(position, description, contract, location, posted_at, company_id, requirements_id, role_id) VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), ?, ?, ?)`

	_, err = jp.DB.Exec(stmt5, jobPostData.Position, jobPostData.Description, jobPostData.Contract, jobPostData.Location, companyID, lastReqInsertedID, lastRoleInsertedID)
	if err != nil {
		return err
	}

	return nil
}

func (jp *JobPostModel) DeleteJobPost(jobPostID int) error {

	stmt := `DELETE FROM jobposts WHERE job_post_id = ?`

	_, err := jp.DB.Exec(stmt, jobPostID)
	if err != nil {
		return err
	}

	return nil
}

func (jp *JobPostModel) EditJobPost(jobPostData EditJopPostFields) error {

	stmt := `UPDATE jobposts
		SET position = ?, description = ?, contract = ?, location = ?
		WHERE job_post_id = ?
		`
	result, err := jp.DB.Exec(stmt, jobPostData.Position, jobPostData.Description, jobPostData.Contract, jobPostData.Location, jobPostData.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("rows affected for jp: %v", rows)
	fmt.Println()

	stmt2 := `UPDATE requirements
		SET requirements_description = ?, requirements_list = ?
		WHERE req_id = ?`

	reqListJSON, err := json.Marshal(jobPostData.Requirements.Items)
	if err != nil {
		fmt.Println("err", err.Error())
		return ErrCouldNotConvertToJSON
	} else {
		fmt.Println()
		fmt.Printf("JSON value of req list:: %s \n", reqListJSON)
		fmt.Println()
	}

	result, err = jp.DB.Exec(stmt2, jobPostData.Requirements.Content, reqListJSON, jobPostData.Requirements.ReqID)
	if err != nil {
		fmt.Println("err reqs", err.Error())
		fmt.Println("err", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	rows, err = result.RowsAffected()
	if err != nil {
		fmt.Println("err", err.Error())
		return err
	}
	fmt.Println()
	fmt.Printf("rows affected for req: %v", rows)
	fmt.Println()

	stmt3 := `UPDATE roles
	SET role_description = ?, role_list = ?
	WHERE role_id = ?`

	roleListJSON, err := json.Marshal(jobPostData.Role.Items)
	if err != nil {
		fmt.Println("err", err.Error())
		return ErrCouldNotConvertToJSON
	} else {
		fmt.Println()
		fmt.Printf("JSON value of role list:: %s \n", roleListJSON)
		fmt.Println()
	}

	result, err = jp.DB.Exec(stmt3, jobPostData.Role.Content, roleListJSON, jobPostData.Role.RoleID)
	if err != nil {
		fmt.Println("err", err.Error())
		// fmt.Println("err", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	rows, err = result.RowsAffected()
	if err != nil {
		fmt.Println("err", err.Error())
		return err
	}
	fmt.Println()
	fmt.Printf("rows affected for role: %v", rows)
	fmt.Println()

	return nil
}

func (jp *JobPostModel) InsertJobApplication(jobPostID, userID int, cvName, cvPath, coverLetter string) error {

	fmt.Printf("jp id %v: \n", jobPostID)
	fmt.Printf("user id %v: \n", userID)
	fmt.Printf("cv name id %v: \n", cvName)
	fmt.Printf("cv path id %v: \n", cvPath)
	fmt.Printf("coverletter id %v: \n", coverLetter)

	//~~ first query the job_applications table to see if there is already a job application for this jp and user
	var err error
	var result sql.Result
	stmt1 := `SELECT jobpost_id, user_id FROM job_applications WHERE user_id = ?`
	var jobApplication JobApplication
	err = jp.DB.QueryRow(stmt1, userID).Scan(&jobApplication.jobPostID, &jobApplication.userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err.Error())
		} else {
			return err
		}
	}

	if jobApplication.jobPostID == jobPostID && jobApplication.userID == userID {
		fmt.Println()
		fmt.Println("user already has applied to this jobpost")
		fmt.Println()
		return ErrUseAlreadyApplied
	}
	// job_applications
	stmt2 := `INSERT INTO job_applications(jobpost_id, user_id, cv_name, cv_path, applied_at, cover_letter)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	if coverLetter == "" {
		fmt.Println("---- no cover letter ----")
		result, err = jp.DB.Exec(stmt2, jobPostID, userID, cvName, cvPath, "none")
	} else {
		fmt.Println("---- with cover letter ----")
		result, err = jp.DB.Exec(stmt2, jobPostID, userID, cvName, cvPath, coverLetter)
	}

	if err != nil {
		return err
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return err
	}

	fmt.Printf("Last insert ID: %d\n", lastInsertID)

	return nil
}
