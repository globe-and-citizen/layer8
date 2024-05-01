describe('My Vue.js App', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5173/');
    });
  
    it('should allow user registration', () => {
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
  
      cy.get('button:contains("Register")').click();
      
      cy.wait('@registerRequest').then((xhr) => {
        expect(xhr.response.statusCode).to.equal(200);
        cy.contains('Registration successful!').should('be.visible');
      });
    });
  
    it('should allow user login', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.intercept('POST', 'http://localhost:5173/api/login').as('loginRequest');
  
      cy.get('button:contains("Login")').click();
  
      cy.wait('@loginRequest').then((xhr) => {
        expect(xhr.response.statusCode).to.equal(200);
        cy.contains('Login successful!').should('be.visible');
      });
    });
  
    it('should allow user to upload profile picture', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.get('button:contains("Login")').click();
  
      cy.wait('@loginRequest').then(() => {
        cy.fixture('profile.jpg').then((fileContent) => {
          cy.get('input[type="file"]').attachFile({
            fileContent: fileContent.toString(),
            fileName: 'profile.jpg',
            mimeType: 'image/jpeg'
          });
        });
        cy.intercept('POST', 'http://localhost:5173/api/profile/upload').as('uploadProfilePictureRequest');
  
        cy.get('button:contains("Upload Profile Picture")').click();
  
        cy.wait('@uploadProfilePictureRequest').then((xhr) => {
          expect(xhr.response.statusCode).to.equal(200);
          cy.contains('Profile picture uploaded successfully!').should('be.visible');
        });
      });
    });
  
    it('should allow user logout', () => {
      cy.get('input[placeholder="default user: tester"]').type('tester');
      cy.get('input[placeholder="default pass: 1234"]').type('1234');
      cy.get('button:contains("Login")').click();
  
      cy.wait('@loginRequest').then(() => {
        cy.get('button:contains("Logout")').click();
  
        cy.contains('You are now logged out!').should('be.visible');
      });
    });
  });
  