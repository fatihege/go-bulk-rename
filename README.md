# usage of bulk-rename
	-append
		specify whether to append or replace the file extension to be renamed
	-ext
		if you want to rename file extensions then use this option
	-files string`
		files to be renamed
	-keep-ext
		if you are renaming filenames choose whether to keep the file extension
	-new string
		new extension name

## examples
	$ bulk-rename -files "./*" -new "new-name"
	$ bulk-rename -files "./*" -new "new-name" -keep-ext
	$ bulk-rename -files "./file1,./file2" -new "new-name"

	$ bulk-rename -ext -files "./*" -new "go"
	$ bulk-rename -ext -files "./*" -new "go" -append
	$ bulk-rename -ext -files "./file1.cpp,./file2.cpp" -new "go"
