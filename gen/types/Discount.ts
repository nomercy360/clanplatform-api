export type Discount = {
    /**
     * @type string | undefined
    */
    code?: string;
    /**
     * @type string | undefined
    */
    created_at?: string;
    /**
     * @type string | undefined
    */
    deleted_at?: string;
    /**
     * @type string
    */
    ends_at: string | null;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type boolean | undefined
    */
    is_active?: boolean;
    /**
     * @type string | undefined
    */
    starts_at?: string;
    /**
     * @type string | undefined
    */
    type?: string;
    /**
     * @type string | undefined
    */
    updated_at?: string;
    /**
     * @type integer | undefined
    */
    usage_count?: number;
    /**
     * @type integer | undefined
    */
    usage_limit?: number;
    /**
     * @type integer | undefined
    */
    value?: number;
};