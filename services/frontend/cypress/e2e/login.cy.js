describe('Login test', () => {
  it('attempts to login', () => {
    cy.visit('http://localhost/signIn');
    cy.get('input[type="text"]').type('a'); // имя тесторого пользователя
    cy.get('input[type="password"]').type('a'); // пароль
    cy.get('button').click();
    cy.location('pathname').should('eq', '/');
    cy.contains('Мероприятия').should('exist');
  })
})
