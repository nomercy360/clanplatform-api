import type { User } from "./User";

 /**
 * @description OK
*/
export type GetAdminUsers200 = User[];

 /**
 * @description OK
*/
export type GetAdminUsersQueryResponse = User[];

 export type GetAdminUsersQuery = {
    Response: GetAdminUsersQueryResponse;
};