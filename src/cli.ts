const { build } = require('gluegun')

async function run(argv) {
  const cli = build()
    .brand('go-dito')
    .src(__dirname)
    .plugins('./node_modules', { matching: 'go-dito-*', hidden: true })
    .help()
    .version()
    .defaultCommand()
    .create()

  const toolbox = await cli.run(argv)

  return toolbox
}

module.exports = { run }
