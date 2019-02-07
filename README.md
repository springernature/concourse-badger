# Concourse Badger

<a href="https://concourse.halfpipe.io/teams/engineering-enablement/pipelines/concourse-badger"><img src="https://badger.halfpipe.io/engineering-enablement/concourse-badger" title="badge"></a>

A little service that renders a pipeline status badge.

It talks directly to the Concourse database to avoid having to deal with atc auth.


## Usage

`GET /<team-name>/<pipeline-name>`


## Development

`./build` and `./run` do what you might expect.



