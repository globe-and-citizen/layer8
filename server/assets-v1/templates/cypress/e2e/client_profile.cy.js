describe('Authentication Page', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5001/client-profile');
    });
  
    it('displays user details', () => {
      cy.get('#app').within(() => {
        cy.contains('Welcome “{{ user.name }}!” Client Portal').should('be.visible');
        cy.contains('Your data').should('be.visible');
        cy.contains('Name:').next().should('have.value', '{{ user.name }}');
        cy.contains('Redirect URI:').next().should('have.value', '{{ user.redirect_uri }}');
        cy.contains('UUID:').next().should('have.value', '{{ user.id }}');
        cy.contains('Secret:').next().should('have.value', '{{ user.secret }}');
      });
    });
  
    it('allows user to copy UUID and Secret to clipboard', () => {
      cy.get('#app').within(() => {
        cy.contains('UUID:').next().within(() => {
          cy.get('button').click();
          cy.get('input').invoke('prop', 'readonly').should('be.true');
          cy.get('input').invoke('val').then((copiedValue) => {
            expect(copiedValue).to.eq('{{ user.id }}');
          });
        });
  
        cy.contains('Secret:').next().within(() => {
          cy.get('button').click();
          cy.get('input').invoke('prop', 'readonly').should('be.true');
          cy.get('input').invoke('val').then((copiedValue) => {
            expect(copiedValue).to.eq('{{ user.secret }}');
          });
        });
      });
    });
  
  });
  