package notionapi

import (
	"encoding/json"
	"time"
)

const (
	// BlockPage is a notion Page
	BlockPage = "page"
	// BlockText is a text block
	BlockText = "text"
	// BlockBookmark is a bookmark block
	BlockBookmark = "bookmark"
	// BlockBulletedList is a bulleted list block
	BlockBulletedList = "bulleted_list"
	// BlockNumberedList is a numbered list block
	BlockNumberedList = "numbered_list"
	// BlockToggle is a toggle block
	BlockToggle = "toggle"
	// BlockTodo is a todo block
	BlockTodo = "to_do"
	// BlockDivider is a divider block
	BlockDivider = "divider"
	// BlockImage is an image block
	BlockImage = "image"
	// BlockHeader is a header block
	BlockHeader = "header"
	// BlockSubHeader is a header block
	BlockSubHeader = "sub_header"
	// BlockSubSubHeader
	BlockSubSubHeader = "sub_sub_header"
	// BlockQuote is a quote block
	BlockQuote = "quote"
	// BlockComment is a comment block
	BlockComment = "comment"
	// BlockCode is a code block
	BlockCode = "code"
	// BlockColumnList is for multi-column. Number of columns is
	// number of content blocks of type TypeColumn
	BlockColumnList = "column_list"
	// BlockColumn is a child of TypeColumnList
	BlockColumn = "column"
	// BlockTable is a table block
	BlockTable = "table"
	// BlockCollectionView is a collection view block
	BlockCollectionView = "collection_view"
	// BlockVideo is youtube video embed
	BlockVideo = "video"
	// BlockFile is an embedded file
	BlockFile = "file"
	// BlockPDF is an embedded pdf file
	BlockPDF = "pdf"
	// BlockGist is embedded gist block
	BlockGist = "gist"
	// BlockTweet is embedded gist block
	BlockTweet = "tweet"
	// BlockEmbed is a generic oembed link
	BlockEmbed = "embed"
	// BlockCallout is a callout
	BlockCallout = "callout"
	// BlockTableOfContents is table of contents
	BlockTableOfContents = "table_of_contents"
)

// BlockPageType defines a type of BlockPage block
type BlockPageType int

const (
	// BlockPageTopLevel is top-level block for the whole page
	BlockPageTopLevel BlockPageType = iota
	// BlockPageSubPage is a sub-page
	BlockPageSubPage
	// BlockPageLink a link to a page
	BlockPageLink
)

// Block describes a block
type Block struct {
	// a unique ID of the block
	ID string `json:"id"`
	// values that come from JSON
	// if false, the page is deleted
	Alive bool `json:"alive"`
	// List of block ids for that make up content of this block
	// Use Content to get corresponding block (they are in the same order)
	ContentIDs   []string `json:"content,omitempty"`
	CopiedFrom   string   `json:"copied_from,omitempty"`
	CollectionID string   `json:"collection_id,omitempty"` // for BlockCollectionView
	// ID of the user who created this block
	CreatedBy   string `json:"created_by"`
	CreatedTime int64  `json:"created_time"`
	// List of block ids with discussion content
	DiscussionIDs []string `json:"discussion,omitempty"`
	// those ids seem to map to storage in s3
	// https://s3-us-west-2.amazonaws.com/secure.notion-static.com/${id}/${name}
	FileIDs   []string        `json:"file_ids,omitempty"`
	FormatRaw json.RawMessage `json:"format,omitempty"`

	// TODO: don't know what this means
	IgnoreBlockCount bool `json:"ignore_block_count,omitempty"`

	// ID of the user who last edited this block
	LastEditedBy   string `json:"last_edited_by"`
	LastEditedTime int64  `json:"last_edited_time"`
	// ID of parent Block
	ParentID    string `json:"parent_id"`
	ParentTable string `json:"parent_table"`
	// not always available
	Permissions *[]Permission          `json:"permissions,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	// type of the block e.g. TypeText, TypePage etc.
	Type string `json:"type"`
	// blocks are versioned
	Version int64    `json:"version"`
	ViewIDs []string `json:"view_ids,omitempty"`

	// Values calculated by us

	// Parent of this block
	Parent *Block `json:"-"`

	// maps ContentIDs array
	Content []*Block `json:"content_resolved,omitempty"`
	// this is for some types like TypePage, TypeText, TypeHeader etc.
	InlineContent []*InlineBlock `json:"inline_content,omitempty"`

	// for BlockPage
	Title string `json:"title,omitempty"`
	// TODO: TitleFull should be Title and we should have
	// GetTitleSimple() function which returns flattened string
	TitleFull []*InlineBlock `json:"title_full,omitempty"`

	// For BlockTodo, a checked state
	IsChecked bool `json:"is_checked,omitempty"`

	// for BlockBookmark
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`

	// for BlockBookmark it's the url of the page
	// for BlockGist it's the url for the gist
	// fot BlockImage it's url of the image, but use ImageURL instead
	// because Source is sometimes not accessible
	// for BlockFile it's url of the file
	// for BlockEmbed it's url of the embed
	Source string `json:"source,omitempty"`

	// for BlockFile
	FileSize string `json:"file_size,omitempty"`

	// for BlockImage it's an URL built from Source that is always accessible
	ImageURL string `json:"image_url,omitempty"`

	// for BlockCode
	Code         string `json:"code,omitempty"`
	CodeLanguage string `json:"code_language,omitempty"`

	// for BlockCollectionView
	// It looks like the info about which view is selected is stored in browser
	CollectionViews []*CollectionViewInfo `json:"collection_views,omitempty"`

	FormatPage     *FormatPage     `json:"format_page,omitempty"`
	FormatBookmark *FormatBookmark `json:"format_bookmark,omitempty"`
	FormatImage    *FormatImage    `json:"format_image,omitempty"`
	FormatColumn   *FormatColumn   `json:"format_column,omitempty"`
	FormatText     *FormatText     `json:"format_text,omitempty"`
	FormatTable    *FormatTable    `json:"format_table,omitempty"`
	FormatVideo    *FormatVideo    `json:"format_video,omitempty"`
	FormatEmbed    *FormatEmbed    `json:"format_embed,omitempty"`
	FormatToggle   *FormatToggle   `json:"format_toggle,omitempty"`
}

// CollectionViewInfo describes a particular view of the collection
type CollectionViewInfo struct {
	CollectionView *CollectionView
	Collection     *Collection
	CollectionRows []*Block
}

// CreatedOn return the time the page was created
func (b *Block) CreatedOn() time.Time {
	return time.Unix(b.CreatedTime/1000, 0)
}

// UpdatedOn returns the time the page was last updated
func (b *Block) UpdatedOn() time.Time {
	return time.Unix(b.LastEditedTime/1000, 0)
}

// IsLinkToPage returns true if block element is a link to a page
// (as opposed to embedded page)
func (b *Block) IsLinkToPage() bool {
	if b.Type != BlockPage {
		return false
	}
	return b.ParentTable == TableSpace
}

// GetPageType returns type of this page
func (b *Block) GetPageType() BlockPageType {
	if b.Parent == nil {
		return BlockPageTopLevel
	}
	if b.ParentID == b.Parent.ID {
		return BlockPageSubPage
	}
	return BlockPageLink
}

// IsPage returns true if block represents a page (either a
// sub-page or a link to a page)
func (b *Block) IsPage() bool {
	return b.Type == BlockPage
}

// IsImage returns true if block represents an image
func (b *Block) IsImage() bool {
	return b.Type == BlockImage
}

// IsCode returns true if block represents a code block
func (b *Block) IsCode() bool {
	return b.Type == BlockCode
}

// FormatToggle describes format for BlockToggle
type FormatToggle struct {
	BlockColor string `json:"block_color"`
}

// FormatPage describes format for BlockPage
type FormatPage struct {
	// /images/page-cover/gradients_11.jpg
	PageCover string `json:"page_cover"`
	// e.g. 0.6
	PageCoverPosition float64 `json:"page_cover_position"`
	PageFont          string  `json:"page_font"`
	PageFullWidth     bool    `json:"page_full_width"`
	// it's url like https://s3-us-west-2.amazonaws.com/secure.notion-static.com/8b3930e3-9dfe-4ba7-a845-a8ff69154f2a/favicon-256.png
	// or emoji like "✉️"
	PageIcon      string `json:"page_icon"`
	PageSmallText bool   `json:"page_small_text"`
	BlockColor    string `json:"block_color"`

	// calculated by us
	PageCoverURL string `json:"page_cover_url,omitempty"`
}

// FormatBookmark describes format for BlockBookmark
type FormatBookmark struct {
	BookmarkIcon string `json:"bookmark_icon"`
}

// FormatImage describes format for BlockImage
type FormatImage struct {
	// comes from notion API
	BlockAspectRatio   float64 `json:"block_aspect_ratio"`
	BlockFullWidth     bool    `json:"block_full_width"`
	BlockPageWidth     bool    `json:"block_page_width"`
	BlockPreserveScale bool    `json:"block_preserve_scale"`
	BlockWidth         float64 `json:"block_width"`
	DisplaySource      string  `json:"display_source,omitempty"`

	// calculated by us
	ImageURL string `json:"image_url,omitempty"`
}

// FormatVideo describes fromat form BlockVideo
type FormatVideo struct {
	BlockWidth         int64   `json:"block_width"`
	BlockHeight        int64   `json:"block_height"`
	DisplaySource      string  `json:"display_source"`
	BlockFullWidth     bool    `json:"block_full_width"`
	BlockPageWidth     bool    `json:"block_page_width"`
	BlockAspectRatio   float64 `json:"block_aspect_ratio"`
	BlockPreserveScale bool    `json:"block_preserve_scale"`
}

// FormatText describes format for BlockText
// TODO: possibly more?
type FormatText struct {
	BlockColor string `json:"block_color,omitempty"`
}

// FormatTable describes format for BlockTable
type FormatTable struct {
	TableWrap       bool             `json:"table_wrap"`
	TableProperties []*TableProperty `json:"table_properties"`
}

// TableProperty describes property of a table
type TableProperty struct {
	Width    int    `json:"width"`
	Visible  bool   `json:"visible"`
	Property string `json:"property"`
}

// FormatColumn describes format for BlockColumn
type FormatColumn struct {
	ColumnRation float64 `json:"column_ratio"` // e.g. 0.5 for half-sized column
}

// FormatEmbed describes format for BlockEmbed
type FormatEmbed struct {
	BlockFullWidth     bool    `json:"block_full_width"`
	BlockHeight        float64 `json:"block_height"`
	BlockPageWidth     bool    `json:"block_page_width"`
	BlockPreserveScale bool    `json:"block_preserve_scale"`
	DisplaySource      string  `json:"display_source"`
}

// Permission describes user permissions
type Permission struct {
	Role   string  `json:"role"`
	Type   string  `json:"type"`
	UserID *string `json:"user_id,omitempty"`
}
