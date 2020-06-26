# Checks a CSS file. Requires use of the 
# Nu HTML checker, which only checks HTML files.
# So you create a stub HTML file, embed the
# contents of the input file into it
# in the form of a style tag, and
# check the resulting file. So the line
# numbers refer to the HTML file, not the
# original CSS file.
OUTFILE="VNU_CHECK.HTML"
# Read the file named on the command line, which is $1, 
# into the variable $INFILE.
INFILE=$(<$1)
# Create the variable $CONTENTS from this HTML, plus
# the contents of $1
read -r -d '' CONTENTS << EOM
<!DOCTYPE html>
<html lang="">
<head>
<title>$1</title>
<style>
$INFILE
</style>
</head>
<body>
<p></p>
</body>
</html>
EOM
## Copy all of this to the output file
echo "$CONTENTS" > $OUTFILE
vnu $OUTFILE

