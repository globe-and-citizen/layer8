const path = require('path');
require('dotenv').config({ path: path.resolve(__dirname, '../../server/.env') });

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
          const dbConn = new Client({
            user: String(process.env.DB_USER),
            password: String(process.env.DB_PASSWORD),
            host: String(process.env.DB_HOST),
            database: String(process.env.DB_NAME),
            port: String(process.env.DB_PORT),
            ssl: process.env.SSL_MODE === 'true' ? { rejectUnauthorized: false } : false,
          });

          try {
            await dbConn.connect();

            // Retrieve user ID by username from tableName
            const userResult = await dbConn.query(`SELECT id FROM ${tableName} WHERE username = $1`, [username]);
            if (userResult.rows.length === 0) {
              throw new Error(`User with username ${username} not found in ${tableName}`);
            }
            const userId = userResult.rows[0].id;

            if (tableName === 'users') {
                // Delete data from user_metadata table using user ID
                await dbConn.query(`DELETE FROM user_metadata WHERE user_id = $1`, [userId]);
                await dbConn.query(`DELETE FROM users WHERE username = $1`, [username]);
            } else {
                await dbConn.query(`DELETE FROM client_traffic_statistics WHERE client_id = $1`, [userId]);
                await dbConn.query(`DELETE FROM clients WHERE username = $1`, [username]);
            }

            await dbConn.end();
            return true;
          } catch (err) {
            await dbConn.end();
            throw new Error(`Database error: ${err.message}`);
          }
        }
      });
    },
  },
});
