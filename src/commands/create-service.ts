import { GluegunToolbox } from 'gluegun';

const buildProjectFiles = async ({
  parameters,
  template: { generate },
  print: { spin, checkmark },
}) => {
  const spinner = spin('Creating project files...')

  const name = parameters.first

  let { options: { module } } = parameters
  if (!module) {
    module = name
  }

  await Promise.all(
    [
      generate({
        template: 'cmd/debug/main.go.ejs',
        target: `${name}/cmd/debug/main.go`,
        props: { module }
      }),
      generate({
        template: 'application/contracts.go.ejs',
        target: `${name}/application/contracts.go`,
        props: { module }
      }),
      generate({
        template: 'domain/contracts.go.ejs',
        target: `${name}/domain/contracts.go`
      }),
      generate({
        template: 'domain/entities.go.ejs',
        target: `${name}/domain/entities.go`
      }),
      generate({
        template: 'domain/value_objects.go.ejs',
        target: `${name}/domain/value_objects.go`
      }),
      generate({
        template: 'infra/contracts.go.ejs',
        target: `${name}/infra/contracts.go`
      }),
      generate({
        template: 'infra/value_objects.go.ejs',
        target: `${name}/infra/value_objects.go`
      }),
      generate({
        template: 'infra/log/log.go.ejs',
        target: `${name}/infra/log/log.go`,
        props: { module }
      }),
      generate({
        template: 'infra/errors/errors.go.ejs',
        target: `${name}/infra/errors/errors.go`,
        props: { module }
      }),
      generate({
        template: 'infra/errors/errors_test.go.ejs',
        target: `${name}/infra/errors/errors_test.go`,
        props: { module }
      }),
      generate({
        template: 'makefile.ejs',
        target: `${name}/Makefile`
      }),
      generate({
        template: 'gitignore.ejs',
        target: `${name}/.gitignore`
      }),
      generate({
        template: 'readme.md.ejs',
        target: `${name}/README.md`,
        props: { name }
      }),
      generate({
        template: 'go.mod.ejs',
        target: `${name}/go.mod`,
        props: { module: module }
      }),
      generate({
        template: '.env-sample.ejs',
        target: `${name}/.env-sample`
      }),
      generate({
        template: '.env-sample.ejs',
        target: `${name}/.env`
      }),
      generate({
        template: 'github/pull_request.md.ejs',
        target: `${name}/.github/PULL_REQUEST_TEMPLATE.md`
      }),
      generate({
        template: 'github/issue_template/bug_report.md.ejs',
        target: `${name}/.github/ISSUE_TEMPLATE/bug_report.md`
      }),
      generate({
        template: 'github/issue_template/feature_request.md.ejs',
        target: `${name}/.github/ISSUE_TEMPLATE/feature_request.md`
      }),
    ]
  )

  spinner.stopAndPersist({ symbol: checkmark, text: 'Project files created' })
}

const downloadDependencies = async ({ parameters, system, print: { spin, checkmark } }) => {
  let spinner = spin('Downloading dependencies...')
  await system.run(`cd ${parameters.first} && make deps`)
  spinner.stopAndPersist({ symbol: checkmark, text: 'Dependencies downloaded' })
}

const createMocks = async ({ parameters, system, print: { spin, checkmark } }) => {
  let spinner = spin('Creating mocks...')
  await system.run(`cd ${parameters.first} && make mock`)
  spinner.stopAndPersist({ symbol: checkmark, text: 'Mocks created' })
}

const runTests = async ({ parameters, system, print: { spin, checkmark } }) => {
  let spinner = spin('Running tests...')
  await system.run(`cd ${parameters.first} && make test`)
  spinner.stopAndPersist({ symbol: checkmark, text: 'Tests passed' })
}

const createGitRepository = async ({ parameters, system, print: { spin, checkmark } }) => {
  let spinner = spin('Creating Git repository...')
  await system.run(`
    cd ${parameters.first} && \
    git init && \
    git add . && \
    git commit -m "First commit"
  `)
  spinner.stopAndPersist({ symbol: checkmark, text: 'Git repository created' })
}

module.exports = {
  name: 'create-service',
  run: async (toolbox: GluegunToolbox) => {
    const {
      parameters,
      print: { error },
    } = toolbox

    if (!parameters.first) {
      error(`Missing service name parameter`)
      process.exit(0)
    }

    await buildProjectFiles(toolbox)
    await downloadDependencies(toolbox)
    await createMocks(toolbox)
    await runTests(toolbox)
    await createGitRepository(toolbox)
  }
}
