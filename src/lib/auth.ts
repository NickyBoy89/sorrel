import { keycloakClientId, keycloakRealm, keycloakUrl } from '../constants';
import Keycloak, { type KeycloakInitOptions } from 'keycloak-js';
import { bearerToken } from '../routes/(app)/stores';

let instance = {
    url: keycloakUrl,
    realm: keycloakRealm,
    clientId: keycloakClientId
};

const initOptions: KeycloakInitOptions = {
    onLoad: "login-required",
};

const keycloak = new Keycloak(instance);

export const initKeycloak = () => {
    if (!keycloak.authenticated) {
        keycloak.init(initOptions)
        .then((auth) => {
            if (auth) {
                console.log("Auth successful!");
            }
            bearerToken.set(keycloak.token);		
        })
        .catch((error) => { 
            console.error("Auth not successful")
            console.error(error)
        });
    }
};