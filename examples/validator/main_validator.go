package validator

import . "github.com/geekymedic/neon/utils/validator"

// define new struct
type BookValidator struct {
    Price *Field
    Name *Field
    ExtraDescPublishingHouse *Field
    ExtraDescAuthor *Field
}

func NewBookValidator(book *Book) *BookValidator {
    return &BookValidator{
        Price: &Field{
            Tag: "Book.Price",
            Value: book.Price,
        }, Name: &Field{
            Tag: "Book.Name",
            Value: book.Name,
        }, ExtraDescPublishingHouse: &Field{
            Tag: "Book.ExtraDesc.PublishingHouse",
            Value: book.ExtraDesc.PublishingHouse,
        }, ExtraDescAuthor: &Field{
            Tag: "Book.ExtraDesc.Author",
            Value: book.ExtraDesc.Author,
        }, 
    }
}