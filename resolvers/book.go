package resolvers

import (
	"context"
	"fmt"
	pb "github.com/les-cours/user-service/protobuf/book"
	"github.com/les-cours/user-service/utils"
)

func (s *Server) GetBook(ctx context.Context, input *pb.BookId) (*pb.BookResponse, error) {
	var book pb.Book
	var err error

	err = s.DB.QueryRow(`
    SELECT id, title, author FROM books WHERE id = $1
    `, input.Id).Scan(&book.Id, &book.Title, &book.Author)

	if err != nil {
		return nil, err
	}

	return &pb.BookResponse{
		Book: &book,
	}, nil
}

func (s *Server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.BookId, error) {

	//GENERATE ID
	bookId := utils.GenerateUUID()
	//INSERT
	_, err := s.DB.Exec(`
INSERT 
INTO
    books
	VALUES 
	    ($1, $2, $3, $4, $5, $6);
`, bookId, in.Title, in.Author, in.Description, in.Pages, in.IsPublished)
	if err != nil {
		return nil, err
	}

	//RETURN

	return &pb.BookId{
		Id: bookId,
	}, nil
}

func (s *Server) UpdateBook(ctx context.Context, in *pb.Book) (*pb.BookId, error) {

	fmt.Println(in.IsPublished)
	fmt.Println(in.Pages)
	if len(in.Title) != 0 {
		err := updateColumn(s, "title", in.Title, in.Id)
		if err != nil {
			return nil, err
		}
	}

	if len(in.Author) != 0 {
		err := updateColumn(s, "author", in.Author, in.Id)
		if err != nil {
			return nil, err
		}
	}

	if len(in.Description) != 0 {
		err := updateColumn(s, "description", in.Description, in.Id)
		if err != nil {
			return nil, err
		}
	}

	if in.Pages != 0 {
		err := updateColumn(s, "pages", in.Pages, in.Id)
		if err != nil {
			return nil, err
		}
	}

	//RETURN

	return &pb.BookId{
		Id: in.Id,
	}, nil
}

func updateColumn(s *Server, column string, value any, id string) error {

	Query := `
	UPDATE public.books
		SET ` + column + `= $1
		WHERE id = $2;
`
	fmt.Println(Query)

	_, err := s.DB.Exec(Query, value, id)
	if err != nil {
		fmt.Sprintf("err in update svx %d", err)
		return err
	}
	return nil

}

func (s *Server) DeleteBook(ctx context.Context, in *pb.BookId) (*pb.BookId, error) {

	_, err := s.DB.Exec(`
DELETE 
FROM 
    books
	WHERE 
	    id = $1;
`, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.BookId{
		Id: in.Id,
	}, nil
}
func (s *Server) GetBooks(ctx context.Context, in *pb.BookPagination) (*pb.BooksResponse, error) {
	var books []*pb.Book

	pagination := pb.Pagination{
		CurrentPage: in.Pagination.CurrentPage,
		PerPage:     in.Pagination.PerPage,
	}

	var ids string
	for i, id := range in.Ids {
		if i == 0 {
			ids = fmt.Sprintf("'%s'", id)
			continue
		}
		ids = fmt.Sprintf("'%s',%s", id, ids)
	}

	QUERY := `
    SELECT 
        id,title, author,description,pages,is_published 
    FROM 
        books 
    WHERE
        id IN(` + ids +
		`) offset $1 row fetch next $2 rows only;`

	rows, err := s.DB.Query(QUERY, (pagination.CurrentPage-1)*pagination.PerPage, pagination.PerPage)

	if err != nil {
		fmt.Sprintf("err in getbooks svx %d", err)
		return nil, err
	}

	fmt.Println(QUERY, (pagination.CurrentPage-1)*pagination.PerPage, pagination.PerPage)
	for rows.Next() {
		var book pb.Book
		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Description, &book.Pages, &book.IsPublished)
		if err != nil {
			fmt.Sprintf("line 130 %d", err)
			return nil, err
		}

		books = append(books, &book)

	}

	fmt.Println(books)
	return &pb.BooksResponse{
		Books: books,
	}, nil
}
