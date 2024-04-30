// describe('Register Client Page', () => {
//   beforeEach(() => {
//     cy.visit('http://localhost:5001/client-register-page')
//   })

//   it('displays the registration form', () => {
//     cy.get('h1').should('contain', 'Register your product')
//     cy.get('input[id="name"]').should('exist')
//     cy.get('input[id="redirect_uri"]').should('exist')
//     cy.get('input[id="username"]').should('exist')
//     cy.get('input[id="password"]').should('exist')
//     cy.get('button').should('contain', 'Register')
//     cy.contains('Already have an account?').should('exist')
//   })

//   it('allows clients to register with valid data', () => {
//     cy.get('input[id="name"]').type('Test Project')
//     cy.get('input[id="redirect_uri"]').type('https://example.com/callback')
//     cy.get('input[id="username"]').type('testuser')
//     cy.get('input[id="password"]').type('password123')
//     cy.get('button').click()
//     cy.url().should('include', 'http://localhost:5001/client-login-page')
//   })

//   it('displays an error message for incomplete registration data', () => {
//     cy.get('button').click()
//     cy.on('window:alert', (message) => {
//       expect(message).to.equal('Please enter all fields!')
//     })
//   })
// })

describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5001/client-login-page')
  })

  it('displays the login form', () => {
    cy.get('h1').should('contain', 'Login')
    cy.get('input#username').should('exist')
    cy.get('input#password').should('exist')
    cy.get('button').should('contain', 'Login')
  })

  it('allows users to login with valid credentials', () => {
    cy.get('input#username').type('hydrolife')
    cy.get('input#password').type('1234')
    cy.get('button').click()
    cy.url().should('include', 'http://localhost:5001/client-profile')
    cy.get('#app').within(() => {
      cy.contains('Welcome “hydrolife!” Client Portal').should('be.visible');
      cy.contains('Your data').should('be.visible');
      cy.contains('Name:').next().invoke('val').should('eq', 'hydrolife');
      cy.contains('Redirect URI:').next().should('have.value', 'hydrolife.com');
      cy.contains('UUID:').next().should('have.value', 'bd2422b6-2357-4f8f-ba46-c1e70c5f0173');
      cy.contains('Secret:').next().should('have.value', 'b333a024c425f1b250e9cd8084093220edbddc7f727ab31797232e48a3d57a59');
    });
    cy.get('#app').within(() => {
      cy.contains('UUID:').next().within(() => {
        cy.get('button').click();
        cy.get('input').invoke('prop', 'readonly').should('be.true');
        // Replace '{{ user.id }}' with actual UUID or known test data
        // cy.get('input').invoke('val').should('eq', 'actual_uuid');
      });
      cy.contains('Secret:').next().within(() => {
        cy.get('button').click();
        cy.get('input').invoke('prop', 'readonly').should('be.true');
        // Replace '{{ user.secret }}' with actual Secret or known test data
        // cy.get('input').invoke('val').should('eq', 'actual_secret');
      });
    });
  })


  // it('redirects to registration page when "Register" link is clicked', () => {
  //   cy.contains('Register').click()
  //   cy.url().should('include', 'http://localhost:5001/client-register-page')
  // })
})

// describe('Authentication Page', () => {
//   beforeEach(() => {
//     cy.visit('http://localhost:5001/client-profile');
//   });

//   it('displays user details', () => {
//     cy.get('#app').within(() => {
//       cy.contains('Welcome “{{ user.name }}!” Client Portal').should('be.visible');
//       cy.contains('Your data').should('be.visible');
//       cy.contains('Name:').next().should('have.value', '{{ user.name }}');
//       cy.contains('Redirect URI:').next().should('have.value', '{{ user.redirect_uri }}');
//       cy.contains('UUID:').next().should('have.value', '{{ user.id }}');
//       cy.contains('Secret:').next().should('have.value', '{{ user.secret }}');
//     });
//   });

//   it('allows user to copy UUID and Secret to clipboard', () => {
//     cy.get('#app').within(() => {
//       cy.contains('UUID:').next().within(() => {
//         cy.get('button').click();
//         cy.get('input').invoke('prop', 'readonly').should('be.true');
//         cy.get('input').invoke('val').then((copiedValue) => {
//           expect(copiedValue).to.eq('{{ user.id }}');
//         });
//       });

//       cy.contains('Secret:').next().within(() => {
//         cy.get('button').click();
//         cy.get('input').invoke('prop', 'readonly').should('be.true');
//         cy.get('input').invoke('val').then((copiedValue) => {
//           expect(copiedValue).to.eq('{{ user.secret }}');
//         });
//       });
//     });
//   });

// });
