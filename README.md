[![Build Status](https://travis-ci.org/easeway/go-cliview.png?branch=master)](https://travis-ci.org/easeway/go-cliview)

# CLI Views
Print complicated data objects in tree view or table view, with customizable formatting and styling support.

# Usage
```go
import (
	cv "github.com/easeway/go-cliview"
)

func ShowTable(data []interface{}) {
	cv.Table{
		Columns: []cv.Column{
			cv.Column{Title: "ID", Field: "id", Width: 16},
			cv.Column{Title: "Name", Field: "name", MaxWidth: 20},
			cv.Column{Title: "Age", Field: "age", MaxWidth: 3, Align: cv.AlignRight},
			cv.Column{Title: "Home Address", Field: "home", Width: -50},		// 50% width
			cv.Column{Title: "Office Address", Field: "office", Width: -50},	// 50% width
		},
		MaxWidth: 80,
	}.Print(data)
}

func ShowJSON(obj map[string]interface{}) {
	cv.NewTree().Print(obj)
}
```

# Features

- Print complicated object in Tree view
	- Left padding support
	- Customizable indent
	- Formatting values
	- Styling map keys and values
- Print array in Table view
	- Customizable column headers
	- Fixed column width
	- Automatically decide column width using actual data
	- Column width in percentage (specified as negative values)
	- Alignment: left, middle, right
	- Auto-ellipsis text
	- Maximum table width (for percentage widths only)
	- Formatting values
	- Styling headers, cells

# Views Details

## Tree view

### Usage

```go
func ShowInTreeView(data interface{}) {
	cv.NewTree().Print(data)
}
```

or

```go
func ShowInTreeView(data interface{}) {
	tv := &cv.Tree{
		Output: cv.Output{		// Output is optional
			Padding: 10,		// left paddings
			Writer: os.Stderr,	// specify output, default is os.Stdout
			Formatter: func (class string, data interface{},, formatter cv.FormatterFunc) string {
				// custom formatter
				// 'class' can be:
				//    - 'tree:key:path'
				//    - 'tree:val:path'
				// here path is from top of "data", e.g.
				//
				// root:
				//     key1:
				//       - name: Jack
				//       - name: Alice
			    //
				// path for "Jack" is 'tree:val:root/key1/0/name',
				// path for "Alice" is 'tree:val:root/key1/1/name',
				// path for "key1" is 'tree:key:root/key1'
				//
				// if you don't know how to format, simply call chained "formatter"
				// return formatter(class, data, nil)
				//
				// When class is 'tree:key:...', it is used for filtering keys rather than renaming keys
				// (renaming keys should use Styler). If an empty string is returned, the key is skipped.
				...
			},
			Styler: func (class, text string, data interface{}) string {
				// class is similar to formatter
				...
			},
		},
		Indent: cv.DefaultIndent,	// override indent, DefaultIndent is 4
	}
	tv.Print(data)
}
```

## Table view

### Usage

```go
func ShowInTableView(data []map[string]interface{}) {
	tv := &cv.Table{
		Output: cv.Output{		// Output is optional
			Padding: 10,		// left paddings
			Writer: os.Stderr,	// specify output, default is os.Stdout
			Formatter: func (class string, data interface{}, formatter cv.FormatterFunc) string {
				// custom formatter
				// 'class' can be: 'table:row:field'
				...
			},
			Styler: func (class, text string, data interface{}) string {
				// class is similar to formatter, it can also be
				// 'table:head:field'
				...
			},
		},
		Columns: []cv.Column{	// required, define the columns
			cv.Column{
				Title: "Display Title",
				Field: "key to fetch data",
				Width: 10,		// >0 for fixed width
								// <0 used as percentage
								// =0 auto decided from data
				MaxWidth: 10,	// limit the maximum column Width
				Align: cv.AlignLeft,	// this is default
										// can be cv.AlignRight, cv.AlignMiddle
				Fetcher: func(col cv.Column, row map[string]interface{}) interface{} {
					// optional function for fetching cell data with special logic.
					// with this function, the column can be a virtual column which doesn't
					// require the Field in the table data but calculated from other fields
				},
				Formatter: ...,  // same as Output.Formatter but operates on column level
			},
			...
		},
		MaxWidth: 80,			// maximum table width
	}
	tv.Print(data)
}
```

# License
X11/MIT
