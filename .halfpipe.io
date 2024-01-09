team: engineering-enablement
pipeline: concourse-badger

triggers:
- type: timer
  cron: "0 3 * * *"

feature_toggles:
- update-pipeline
- github-statuses

tasks:
- type: docker-compose
  name: build
  save_artifacts: [.]

- type: deploy-cf
  name: deploy
  api: ((cloudfoundry.api-snpaas))
  space: halfpipe
  deploy_artifact: .
  vars:
    POSTGRES_USERNAME: ((halfpipe-concourse-db-pg13.username_read))
    POSTGRES_PASSWORD: ((halfpipe-concourse-db-pg13.password_read))
    POSTGRES_HOST: ((halfpipe-concourse-db-pg13.host))
