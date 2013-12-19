iDoneThis Go API
----------------

This library was written by observing the way that the iDoneThis website functions
and noticing that it uses a JSON API internally to retrieve information. This library
was then written to take advantage of that API so I could write little programs to 
watch iDoneThis, post to iDoneThis, etc. There's a bit of screen scraping to get to 
the JSON api, which is in the NewClient function.

As I could not find API documentation for
the JSON API, I did the best I could as far as finding urls and options. If you have
documentation, or new urls that aren't part of the client, feel free to file an issue
with the information or even write it up as a pull request.

Todo's

- Other done filters?
- Parse all fields of Team struct
- Examples of origin field on dones
- Examples of rawintegrationdata on dones
- Date parser for the Done struct
- Commenting on Dones
- Documenting API

[Documentation](http://godoc.org/github.com/acsellers/idonethis#DoneFilter)
