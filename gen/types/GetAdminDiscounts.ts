import type { Discount } from "./Discount";

 /**
 * @description OK
*/
export type GetAdminDiscounts200 = Discount[];

 /**
 * @description OK
*/
export type GetAdminDiscountsQueryResponse = Discount[];

 export type GetAdminDiscountsQuery = {
    Response: GetAdminDiscountsQueryResponse;
};