const { defineConfig } = require("cypress");
const { Client } = require('pg')

module.exports = defineConfig({
  component: {
    devServer: {
      framework: "vue-cli",
      bundler: "webpack",
    },
  },

  e2e: {
    setupNodeEvents(on, config) {
      on("task", {
        async connectDB(){
          const client = new Client({
            user: "postgres",
            password: "1234",
            host: "localhost",
            database: "ResourceServer",
            ssl: false,
            port: 5432
          })
          await client.connect()
          const username = "testuser"
          const res = await client.query(`DELETE FROM users WHERE username = '${username}'`)
          await client.end()
          return res.rows;
        }
      })
    },
  },
});
