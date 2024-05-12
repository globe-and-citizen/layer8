describe('WGP', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5173/');
    });
  
    it('should allow user registration', () => {
      cy.contains('a.block', "Don't have an account? Register").click();
      cy.get('input[placeholder="Username"]').type('newuser');
      cy.get('input[placeholder="Password"]').type('password123');
      cy.fixture('profile.jpg').then((fileContent) => {
        cy.get('input[type="file"]').attachFile({
          fileContent: fileContent.toString(),
          fileName: 'profile.jpg',
          mimeType: 'image/jpeg'
        });
      });
      cy.intercept('POST', 'http://localhost:5173/api/register').as('registerRequest');
      cy.get('button:contains("Register")').click({force: true});
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Registration successful!')
      })
    });
  
    it('should login with Layer8', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.intercept('POST', 'http://localhost:5173/api/login').as('loginRequest');
  
      cy.get('button:contains("Login")').click();
  
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Login successful!')
      })
    });

    it('should login Anonymously', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.get('button:contains("Login")').click();
    
      // cy.contains('button.btn', 'Login with Layer8').click();
    });

    it('should allow user to upload profile picture', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.get('button:contains("Login")').click();
    
      cy.contains('button.btn', 'Login Anonymously').click();
      cy.url().should('include', 'http://localhost:5173/home');
      for (let i = 0; i < 10; i++) {
        cy.contains('button.btn', 'Get Next Poem').click();
      }
      
      cy.contains('button.btn', 'Logout').click();
    });    
  });
  