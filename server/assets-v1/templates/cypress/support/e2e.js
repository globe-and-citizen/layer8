// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

Cypress.Commands.add('deleteRegisteredUser', () => {
    cy.request({
      method: 'DELETE',
      url: 'http://localhost:5001/api/users',
    }).then((response) => {
      expect(response.status).to.eq(200);
      cy.log('Registered user deleted successfully');
    }).catch((error) => {
      cy.log(`Error deleting registered user: ${error}`);
    });
  });
  