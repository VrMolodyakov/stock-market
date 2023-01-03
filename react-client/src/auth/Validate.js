export const validLogin = new RegExp(
    /^[\w.-]{3,19}[0-9a-zA-Z]$/
 );
 export const validPassword = new RegExp(/^(?=.*?[A-Za-z])(?=.*?[0-9]).{6,}$/);
