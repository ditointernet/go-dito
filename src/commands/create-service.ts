import { GluegunToolbox } from 'gluegun';

module.exports = {
  name: 'create-service',
  run: async (toolbox: GluegunToolbox) => {
    const {
      parameters,
      template: { generate },
      print: { info, error }
    } = toolbox

    const name = parameters.first
    if (!name) {
      return error('😠 Missing service name parameter')
    }

    let { options: { module: moduleName } } = parameters
    if (!moduleName) {
      moduleName = name
    }

    await Promise.all(
      [
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
          props: { module: moduleName }
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

    info(`✔️ Service ${name} created!`)
  }
}
