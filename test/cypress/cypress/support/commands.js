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

import 'cypress-file-upload';

Cypress.Commands.add('loginWGP', () => {
    cy.visit('http://localhost:5173')
    cy.get('input').eq(0).type('testuser')
    cy.get('input').eq(1).type('testpassword')
    cy.get('button').eq(1).click()
    cy.wait(500)
})

Cypress.Commands.add('loginAnonymouslyWGP', () => {
    cy.visit('http://localhost:5173')
    cy.get('input').eq(0).type('testuser')
    cy.get('input').eq(1).type('testpassword')
    cy.get('button').eq(1).click()
    cy.wait(500)
    cy.contains('Login Anonymously').click()
})
