describe('My Vue.js App', () => {
    beforeEach(() => {
      cy.visit('/');
    });
  
    it('should display the navbar', () => {
      cy.get('nav').should('exist');
    });
  
    it('should display the user information after login', () => {
      localStorage.setItem("SP_TOKEN", "dummySpToken");
      localStorage.setItem("L8_TOKEN", "dummyL8Token");
  
      cy.intercept('GET', '**/nextpoem', { fixture: 'nextPoem.json' }).as('getPoem');
  
      cy.wait('@getPoem').then(() => {
        cy.contains('Welcome user123!').should('be.visible');
        cy.get('h4:contains("Username: JohnDoe")').should('be.visible');
        cy.get('h4:contains("Country: USA")').should('be.visible');
        cy.contains('Email Verified: Email is verified!').should('be.visible');
      });
    });
  
    it('should log out the user when the logout button is clicked', () => {
      localStorage.setItem("SP_TOKEN", "dummySpToken");
      localStorage.setItem("L8_TOKEN", "dummyL8Token");
  
      cy.intercept('GET', '**/nextpoem', { fixture: 'nextPoem.json' }).as('getPoem');
  
      cy.wait('@getPoem').then(() => {
        cy.get('button:contains("Logout")').click();
  
        cy.url().should('include', '/loginRegister');
        cy.get('h1').should('contain', 'Login');
        cy.get('input#username').should('exist');
        cy.get('input#password').should('exist');
        cy.get('button:contains("Login")').should('exist');
      });
    });
  });
  