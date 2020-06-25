
OUTFILE=foo.txt
# Read the file named on the command line, which is $1, 
# into the variable $INFILE.
INFILE=$(<$1)
# Create the variable $CONTENTS from this HTML, plus
# the contents of $1
read -r -d '' CONTENTS << EOM
<!DOCTYPE html>
<html lang="">
<head>
<title>$OUTFILE</title>
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


