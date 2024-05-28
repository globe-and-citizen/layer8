describe('WGP', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5173/');
  });

  // it('should allow user registration', () => {
  //   cy.contains('a.block', "Don't have an account? Register").click();
  //   cy.get('input[placeholder="Username"]').type('newuser');
  //   cy.get('input[placeholder="Password"]').type('password123');
  //   cy.fixture('profile.jpg').then((fileContent) => {
  //     cy.get('input[type="file"]').attachFile({
  //       fileContent: fileContent.toString(),
  //       fileName: 'profile.jpg',
  //       mimeType: 'image/jpeg'
  //     });
  //   });
  //   cy.intercept('POST', 'http://localhost:5173/api/register').as('registerRequest');
  //   cy.get('button:contains("Register")').click({force: true});
  //   cy.on('window:alert', (message) => {
  //     expect(message).to.equal('Registration successful!')
  //   })
  // });

  // it('should login with Layer8', () => {
  //   cy.get('input[placeholder="default user: tester"]').type('tester');
  //   cy.get('input[placeholder="default pass: 1234"]').type('1234');
  //   cy.intercept('POST', 'http://localhost:5173/api/login').as('loginRequest');

  //   cy.get('button:contains("Login")').click();

  //   cy.on('window:alert', (message) => {
  //     expect(message).to.equal('Login successful!')
  //   })
  // });

  it('should login Anonymously', () => {
    cy.get('input[placeholder="default user: tester"]').type('tester');
    cy.get('input[placeholder="default pass: 1234"]').type('1234');
    cy.get('button:contains("Login")').click();

    cy.intercept('GET', 'http://localhost:5001/login?next=/authorize?client_id=notanid&scope=read%3Auser', {
      statusCode: 200,
      body: { message: 'Mock response from popup server' },
    }).as('popupRequest');

    // cy.contains('button.btn', 'Login with Layer8').click();


    // Trigger the action that opens the popup
    // cy.get('.open-popup-button').click();

    // Wait for the request to the popup server to complete
    cy.wait('@popupRequest');

    // Verify that the popup funcionality is correctly integrated
    // For example, check that the data from the popup is displayed in the main application
    cy.contains('.popup-data', 'Mock response from popup server').should('exist');

    // // Simulate interaction witth the popup
    // cy.url({ timeout: 10000 }).should('include', 'http://localhost:5001/login?next=/authorize?client_id=notanid&scope=read%3Auser');

    // // Now, you can interact with elements inside the popup
    // cy.get('input[name="username"]').type('your_username');
    // cy.get('input[name="password"]').type('your_password');
    // cy.get('form').submit(); // Submit the form

    // // Wait for the authentication process to complete (you might need to adjust the URL or add a different wait condition)
    // cy.wait(5000); // Adjust the wait time as needed
  });

  // it('should allow user to upload profile picture', () => {
  //   cy.get('input[placeholder="default user: tester"]').type('tester');
  //   cy.get('input[placeholder="default pass: 1234"]').type('1234');
  //   cy.get('button:contains("Login")').click();

  //   cy.contains('button.btn', 'Login Anonymously').click();
  //   cy.url().should('include', 'http://localhost:5173/home');
  //   for (let i = 0; i < 10; i++) {
  //     cy.contains('button.btn', 'Get Next Poem').click();
  //   }

  //   cy.contains('button.btn', 'Logout').click();
  // });
});
