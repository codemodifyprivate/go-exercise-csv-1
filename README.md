# GRAX - Go exercise

This is a coding exercise to show your skills in handling CSV files using the Go language. It's based on real challenges we face at GRAX.

The workflow is simple:

- Clone this repository
- Read the problem statement below
- Spend about 90 minutes writing code/tests for it
- Stick to the standard library
- Feel free to consult any docs about it
- Proceed and commit as you would in a normal project at work

There are obviously no right or wrong answers here; this is just a foundation for us to have a conversation about your motivations for taking a particular approach, the compromises involved, and your thoughts in scaling and operating this kind of tool in production.

So let's dive in!

## A short story on why would anyone pay for Salesforce backups

In case you're not familiar, Salesforce can be described as "Microsoft Access for the Cloud". In some level it really is just a database where enterprises keep their business data.

Different than Access, though, you can't just ask Salesforce to export any table as CSV. The platform has so many constraints they could be a Jeopardy category of its own, going for many, many seasons.

Suffice to say these constraints shape the CSVs exported from Salesforce. At a certain scale your data has to be divided into smaller CSV files:

- When a table has too many rows, it's paginated
- When a table has too many columns, it's fragmented

The goal of this exercise is to write code to consolidate different CSV files representing the same data into a single well formed CSV.

## Program structure and rules

The entrypoint of your program should be in `main.go`: it receives a directory containing CSVs, reads each file, combines and merge them, outputting the final CSV to stdout.

You can run it using the supplied test data, like:

```
go run main.go -dir data/simple
```

A few things to keep in mind about the data you'll work with:

- The first column on every CSV is the "Id", which you can use to match different fragments of the same row
- There's no guarantee the rows will be sorted
- There's no guarantee the files are sorted in any meaningful way, either

And rules about the data you'll produce:

- Please sort the rows by their ID, for convenience
- If there are gaps in the data, fill in with blanks

## Example

Given 4 CSV files:

```
Id,FirstName,LastName
1,Amy,Adams
2,John,Malkovich
```

```
Id,Phone,Email
1,310-111-1111,contact@amyadams.com
2,213-222-2222,john@malkovich.com
```

```
Id,FirstName,LastName
3,Larry,David
4,Michelle,Wolf
```

```
Id,Phone,Email
4,213-444-444,mwolf@comcast.net
```

The consolidated CSV could be:

```
Id,FirstName,LastName,Phone,Email
1,Amy,Adams,310-111-1111,contact@amyadams.com
2,John,Malkovich,213-222-2222,john@malkovich.com
3,Larry,David,,
4,Michelle,Wolf,213-444-444,mwolf@comcast.net
```

Note the columns could be in a different order: `Id,Phone,Email,FirstName,Lastname`.

