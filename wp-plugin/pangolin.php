<?php
/**
 * Plugin Name: Pangolin
 * Description: Provides API key authentication for Pangolin CLI
 * Version: 1.0.0
 * Author: Pangolin
 * License: MIT
 */

if (!defined('ABSPATH')) {
    exit;
}

define('PANGOLIN_VERSION', '1.0.0');
define('PANGOLIN_PLUGIN_DIR', plugin_dir_path(__FILE__));

require_once PANGOLIN_PLUGIN_DIR . 'includes/class-pangolin-api.php';
require_once PANGOLIN_PLUGIN_DIR . 'includes/class-pangolin-settings.php';

class Pangolin {
    private static $instance = null;

    public static function get_instance() {
        if (null === self::$instance) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    private function __construct() {
        add_action('rest_api_init', array($this, 'register_routes'));
        add_action('admin_menu', array($this, 'add_admin_menu'));
        add_action('admin_init', array($this, 'register_settings'));
    }

    public function register_routes() {
        Pangolin_API::register_routes();
    }

    public function add_admin_menu() {
        add_options_page(
            'Pangolin',
            'Pangolin',
            'manage_options',
            'pangolin',
            array('Pangolin_Settings', 'render_page')
        );
    }

    public function register_settings() {
        Pangolin_Settings::register_settings();
    }
}

function pangolin() {
    return Pangolin::get_instance();
}

pangolin();
