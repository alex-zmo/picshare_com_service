package database

import (
	"fmt"
	"strconv"
)


func createCommentsTable() {
	createCommentsStmt := `CREATE TABLE IF NOT EXISTS comment (
		id INT AUTO_INCREMENT,
		user_id INT,
		pic_id INT,
		content VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		reply_of INT DEFAULT 0,
		likes INT DEFAULT 0,
		deleted TINYINT(1) DEFAULT 0,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (pic_id) REFERENCES pic(id)
	);`

	_, err := db.Exec(createCommentsStmt)
	if err != nil {
		fmt.Println(err)
	}
}

func InsertComment(userId, picId int, comment string) string {
	insertCommentsStmt := `INSERT INTO comment (user_id, pic_id, content) VALUES (?, ?, ?);`

	if userId == -1 || picId == -1 {
		return "Error in Insert Comment, Invalid user or pic ID"
	}
	_, err := db.Exec(insertCommentsStmt, userId, picId, comment)
	if err != nil {
		fmt.Println(err)
		return "Error in Insert Comment"
	}
	return ""
}


func GetPicIdFromLink(link string) int {
	retrieveLinkStmt := `SELECT pic_id FROM link WHERE link = ? LIMIT 1
	`
	res, err := db.Query(retrieveLinkStmt, link)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer closeRows(res)

	picId := -1
	if res.Next() {
		err := res.Scan(&picId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return picId
}
func GetUserIdFromPicId(picId int) int {
	retrieveIdStmt := `SELECT user_id FROM pic WHERE id = ? LIMIT 1`
	res, err := db.Query(retrieveIdStmt, picId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer closeRows(res)

	userId := -1
	if res.Next() {
		err := res.Scan(&userId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return userId
}

func ListComments(picId, requestUserId int) [][]string{
  listCommentStmt := `SELECT id, user_id, content, created_at FROM comment WHERE pic_id = ? AND deleted = 0
	`
	result := make([][]string,0)

	deletable := false

	picUserId := GetUserIdFromPicId(picId)
	if picUserId == -1 {
		fmt.Println("Error in get User Id From Pic Id on List Comments")
		return result
	}

	if  picUserId == requestUserId {
		deletable = true
	}

	res, err := db.Query(listCommentStmt, picId)
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer closeRows(res)

	commentId := ""
	userId := ""
	comment := ""
	createdAt := ""


	for res.Next() {
		res.Scan(&commentId, &userId, &comment, &createdAt)
		line := make([]string, 0)

		line = append(line, commentId)
		line = append(line, userId)
		line = append(line, createdAt)
		line = append(line, comment )
		line = append(line, strconv.FormatBool(deletable || strconv.Itoa(requestUserId) == userId) )
		result = append(result, line)
	}
	return result
}

func DeleteComment(commentId int) {
	deleteCommentStmt := `UPDATE comment SET deleted = 1 WHERE id = ?`

	_, err := db.Exec(deleteCommentStmt, commentId)
	if err != nil {
		fmt.Println(err)
	}
}