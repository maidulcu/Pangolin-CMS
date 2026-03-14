<?php

class Pangolin_API {

    public static function register_routes() {
        register_rest_route('pangolin/v1', '/key', array(
            'methods' => 'POST',
            'callback' => array(__CLASS__, 'generate_key'),
            'permission_callback' => array(__CLASS__, 'permissions_check'),
        ));

        register_rest_route('pangolin/v1', '/key', array(
            'methods' => 'DELETE',
            'callback' => array(__CLASS__, 'delete_key'),
            'permission_callback' => array(__CLASS__, 'permissions_check'),
        ));

        register_rest_route('pangolin/v1', '/validate', array(
            'methods' => 'GET',
            'callback' => array(__CLASS__, 'validate_key'),
            'permission_callback' => '__return_true',
        ));

        register_rest_route('pangolin/v1', '/info', array(
            'methods' => 'GET',
            'callback' => array(__CLASS__, 'get_site_info'),
            'permission_callback' => array(__CLASS__, 'authenticate'),
        ));
    }

    public static function permissions_check($request) {
        return current_user_can('edit_posts');
    }

    public static function authenticate($request) {
        $auth_header = $request->get_header('Authorization');
        
        if (!$auth_header) {
            return new WP_Error('rest_notAuthorized', 'No authorization header', array('status' => 401));
        }

        if (strpos($auth_header, 'Bearer ') !== 0) {
            return new WP_Error('rest_notAuthorized', 'Invalid authorization format', array('status' => 401));
        }

        $api_key = substr($auth_header, 7);
        $stored_key = get_option('pangolin_api_key');

        if (!$stored_key || !hash_equals($stored_key, $api_key)) {
            return new WP_Error('rest_notAuthorized', 'Invalid API key', array('status' => 401));
        }

        return true;
    }

    public static function generate_key($request) {
        $api_key = wp_generate_password(32, false);
        
        update_option('pangolin_api_key', $api_key);
        update_option('pangolin_api_created', current_time('mysql'));

        return rest_ensure_response(array(
            'success' => true,
            'api_key' => $api_key,
            'message' => 'API key generated successfully',
        ));
    }

    public static function delete_key($request) {
        delete_option('pangolin_api_key');
        delete_option('pangolin_api_created');

        return rest_ensure_response(array(
            'success' => true,
            'message' => 'API key deleted successfully',
        ));
    }

    public static function validate_key($request) {
        $auth_header = $request->get_header('Authorization');
        
        if (!$auth_header || strpos($auth_header, 'Bearer ') !== 0) {
            return new WP_Error('rest_notAuthorized', 'No API key provided', array('status' => 401));
        }

        $api_key = substr($auth_header, 7);
        $stored_key = get_option('pangolin_api_key');

        if ($stored_key && hash_equals($stored_key, $api_key)) {
            return rest_ensure_response(array(
                'valid' => true,
                'site_url' => get_site_url(),
            ));
        }

        return new WP_Error('rest_notAuthorized', 'Invalid API key', array('status' => 401));
    }

    public static function get_site_info($request) {
        return rest_ensure_response(array(
            'site_url' => get_site_url(),
            'site_name' => get_bloginfo('name'),
            'rest_url' => rest_url(),
            'sitemap_urls' => array(
                get_site_url() . '/sitemap.xml',
                get_site_url() . '/wp-sitemap.xml',
            ),
        ));
    }
}
