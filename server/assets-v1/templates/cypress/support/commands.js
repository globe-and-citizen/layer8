// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })


// Cypress.Commands.add('deleteRegisteredUser', (username, type) => {
//     return cy.request({
//         method: 'DELETE',
//         url: 'http://localhost:5001/api/v1/delete-user',
//         body: {
//             username: username,
//             type: type
//         },
//     })
// });

const { Client } = require('pg');

Cypress.Commands.add('deleteRegisteredUser', (username) => {
    const client = new Client({
        user: 'postgres',
        host: 'localhost',
        database: 'ResourceServer',
        password: '1234',
        port: 5432, // default PostgreSQL port
    });

    return client.connect()
        .then(() => {
            return client.query(`DELETE FROM users WHERE username = '${username}'`);
        })
        .then(() => {
            client.end();
        })
        .catch(err => {
            console.error('Error deleting user from database:', err);
            client.end();
        });
});
