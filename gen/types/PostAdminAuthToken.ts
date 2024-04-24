import { AuthUser } from "./AuthUser";
import type { UserWithToken } from "./UserWithToken";

 /**
 * @description OK
*/
export type PostAdminAuthToken200 = UserWithToken;

 /**
 * @description User data
*/
export type PostAdminAuthTokenMutationRequest = AuthUser;

 /**
 * @description OK
*/
export type PostAdminAuthTokenMutationResponse = UserWithToken;

 export type PostAdminAuthTokenMutation = {
    Response: PostAdminAuthTokenMutationResponse;
    Request: PostAdminAuthTokenMutationRequest;
};