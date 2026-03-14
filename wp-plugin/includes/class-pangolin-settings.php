<?php

class Pangolin_Settings {

    public static function register_settings() {
        register_setting('pangolin', 'pangolin_api_key');
        register_setting('pangolin', 'pangolin_api_created');
    }

    public static function render_page() {
        $api_key = get_option('pangolin_api_key');
        $api_created = get_option('pangolin_api_created');
        
        ?>
        <div class="wrap">
            <h1>Pangolin Settings</h1>
            
            <form method="post" action="">
                <?php settings_fields('pangolin'); ?>
                
                <table class="form-table">
                    <tr>
                        <th scope="row">API Key</th>
                        <td>
                            <?php if ($api_key): ?>
                                <code style="display: block; padding: 10px; background: #f6f7f7; border: 1px solid #c3c4c7; margin-bottom: 10px;">
                                    <?php echo esc_html($api_key); ?>
                                </code>
                                <p class="description">
                                    Created: <?php echo esc_html($api_created); ?>
                                </p>
                                <p>
                                    <button type="submit" name="pangolin_delete_key" class="button button-link-delete">
                                        Delete API Key
                                    </button>
                                </p>
                            <?php else: ?>
                                <p class="description">No API key generated yet.</p>
                            <?php endif; ?>
                        </td>
                    </tr>
                </table>
                
                <p class="submit">
                    <button type="submit" name="pangolin_generate_key" class="button button-primary">
                        <?php echo $api_key ? 'Regenerate API Key' : 'Generate API Key'; ?>
                    </button>
                </p>
            </form>

            <hr>
            
            <h2>Usage with Pangolin CLI</h2>
            <p>After generating an API key, run the following command:</p>
            <pre style="background: #f6f7f7; padding: 15px; border: 1px solid #c3c4c7;">pangolin init -u <?php echo get_site_url(); ?> -k YOUR_API_KEY</pre>
            
            <h3>Available Endpoints</h3>
            <ul>
                <li><code><?php echo rest_url('pangolin/v1/validate'); ?></code> - Validate API key</li>
                <li><code><?php echo rest_url('pangolin/v1/info'); ?></code> - Get site info</li>
            </ul>
        </div>
        <?php

        if (isset($_POST['pangolin_generate_key'])) {
            $api_key = wp_generate_password(32, false);
            update_option('pangolin_api_key', $api_key);
            update_option('pangolin_api_created', current_time('mysql'));
            echo '<div class="notice notice-success"><p>API key generated!</p></div>';
        }

        if (isset($_POST['pangolin_delete_key'])) {
            delete_option('pangolin_api_key');
            delete_option('pangolin_api_created');
            echo '<div class="notice notice-success"><p>API key deleted!</p></div>';
        }
    }
}
