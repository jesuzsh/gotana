# Gotana
> Client to access match data for users in the game, Halo: Infinite.

## Functionality
Each of the following steps are done concurrently.
* Fetch a user's match data, using [AutoCode's Halo
  API](https://autocode.com/halo/)
* Zip data from all the responses
* Persist data in an S3 bucket
