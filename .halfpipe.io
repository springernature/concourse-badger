team: engineering-enablement
pipeline: concourse-badger
cron_trigger: "0 3 * * *"
 
feature_toggles:
- update-pipeline

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
    POSTGRES_USERNAME: ((concourse-db.username_read))
    POSTGRES_PASSWORD: ((concourse-db.password_read))
    POSTGRES_HOST: ((concourse-db.host))
