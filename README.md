# Gotana
> Client to access match data for users in the game, Halo: Infinite.

## Functionality
Each of the following steps are done concurrently.
* Fetch a user's match data, using [AutoCode's Halo
  API](https://autocode.com/halo/)
* Zip data from all the responses
* Persist data in an S3 bucket

## Components
The following are the essential components of the architecture.

### Autocode endpoints
Deployed JavaScript/Python code snippets that interact with an Autocode client 
library. These deployed endpoints act as a light wrapperover Autocode's Halo 
client library, so that they may be leveraged to get match data. Since there is
a max of 25 matches returned per request, N being the total matches and greater
than 25, N/25 is the number of requests needed to get all the match data for a
particular user.

### Concurrency
Goroutines and channels are heavily relied on at every stage of the pipeline.
Starting with perfoming the required number of requests to retrieve data for
all the matches played by a user (Stage 1). Additionally, the gzip package is
used to compress the json responses (Stage 2). Persist compressed json files
to S3 for further analysis (Stage 3). Each stages kicks off numerous goroutines
and the entire pipeline is completed in a matter of seconds (<3).

## Architecture
The following diagram illustrates the current state of the overall pipeline:

![alt text](https://github.com/jesuzsh/gotana/blob/main/architecture.png?raw=true)
