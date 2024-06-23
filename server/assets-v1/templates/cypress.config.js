const { defineConfig } = require("cypress");
const { Client } = require('pg');
require('dotenv').config();

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
        async deleteUser({ username, tableName }) {
          console.log('DB_USER:', String(process.env.DB_USER));
          console.log('DB_PASSWORD:', String(process.env.DB_PASSWORD));
          console.log('DB_HOST:', String(process.env.DB_HOST));
          console.log('DB_NAME:', String(process.env.DB_NAME));
          console.log('DB_PORT:', String(process.env.DB_PORT));

          const client = new Client({
            user: "layer8development",
            password: "cRACtIeRChYsiNFinuMp",
            host: "localhost",
            database: "development",
            port: "5433",
            ssl: false,
          });

          try {
            await client.connect();
            await client.query(`DELETE FROM ${tableName} WHERE username = $1`, [username]);
            await client.end();
            return true;
          } catch (err) {
            await client.end();
            throw new Error(`Database error: ${err.message}`);
          }
        }
      });
    },
  },
});
