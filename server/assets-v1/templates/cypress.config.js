const path = require('path');
require('dotenv').config({ path: path.resolve(__dirname, '../../.env') });

const { defineConfig } = require("cypress");
const { Client } = require('pg');

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
          const client = new Client({
            user: String(process.env.DB_USER),
            password: String(process.env.DB_PASSWORD),
            host: String(process.env.DB_HOST),
            database: String(process.env.DB_NAME),
            port: String(process.env.DB_PORT),
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
