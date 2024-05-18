package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var DefaultCompanyIcon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 512"><path d="M384 96V320H64L64 96H384zM64 32C28.7 32 0 60.7 0 96V320c0 35.3 28.7 64 64 64H181.3l-10.7 32H96c-17.7 0-32 14.3-32 32s14.3 32 32 32H352c17.7 0 32-14.3 32-32s-14.3-32-32-32H277.3l-10.7-32H384c35.3 0 64-28.7 64-64V96c0-35.3-28.7-64-64-64H64zm464 0c-26.5 0-48 21.5-48 48V432c0 26.5 21.5 48 48 48h64c26.5 0 48-21.5 48-48V80c0-26.5-21.5-48-48-48H528zm16 64h32c8.8 0 16 7.2 16 16s-7.2 16-16 16H544c-8.8 0-16-7.2-16-16s7.2-16 16-16zm-16 80c0-8.8 7.2-16 16-16h32c8.8 0 16 7.2 16 16s-7.2 16-16 16H544c-8.8 0-16-7.2-16-16zm32 160a32 32 0 1 1 0 64 32 32 0 1 1 0-64z" fill="#8a8eb4"/></svg>`

var DefaultCompanyBgColor = `hsl(234, 100%, 94%)`

type JobPostModelInterface interface {
	Latest() ([]JobPost, error)
	FilterPosts(position, location, contract string) ([]JobPost, error)
	Get(id int) (JobPost, error)
	InsertCompany(name, logoSVG, logoBgColor, website string) error
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
	Name        string
	LogoSVG     string
	LogoBgColor string
	Website     string
}

type Requirements struct {
	RequirementsDescription string
	RequirementsList        []string
}
type Role struct {
	RoleDescription string
	RoleList        []string
}

type JobPostModel struct {
	DB *sql.DB
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

	stmt := `SELECT jp.job_post_id, jp.position, jp.description, jp.contract, jp.location, jp.posted_at, cp.name, cp.logo_svg, cp.logo_bg_color, cp.website, rq.requirements_description, rq.requirements_list, rl.role_description, rl.role_list
		FROM jobposts AS jp
		INNER JOIN companies AS cp ON jp.company_id = cp.company_id
		INNER JOIN requirements AS rq ON jp.requirements_id = rq.req_id
		INNER JOIN roles AS rl ON jp.role_id = rl.role_id
		WHERE jp.job_post_id = ?`

	row := jp.DB.QueryRow(stmt, id)
	var jobPost JobPost

	var rqList any
	var roleList any

	err := row.Scan(&jobPost.ID, &jobPost.Position, &jobPost.Description, &jobPost.Contract, &jobPost.Location, &jobPost.PostedAt, &jobPost.Company.Name, &jobPost.Company.LogoSVG, &jobPost.Company.LogoBgColor, &jobPost.Company.Website, &jobPost.Requirements.RequirementsDescription, &rqList, &jobPost.Role.RoleDescription, &roleList)

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
