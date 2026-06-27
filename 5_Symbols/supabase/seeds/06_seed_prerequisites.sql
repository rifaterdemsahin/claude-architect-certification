-- Insert sample prerequisites for Module 1
INSERT INTO public.prerequisites (module_number, video_number, install_name, install_command, verification_command)
VALUES 
(1, 1, 'Google Cloud CLI', 'curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-darwin-x86_64.tar.gz && tar -xf google-cloud-cli-darwin-x86_64.tar.gz && ./google-cloud-sdk/install.sh', 'gcloud version'),
(1, 1, 'Docker Desktop', 'brew install --cask docker', 'docker --version'),
(1, 2, 'Fly.io CLI', 'curl -L https://fly.io/install.sh | sh', 'flyctl version'),
(1, 2, 'Supabase CLI', 'brew install supabase/tap/supabase', 'supabase --version'),
(2, 1, 'Node.js', 'nvm install 20', 'node -v && npm -v'),
(2, 2, 'Python 3.11', 'brew install python@3.11', 'python3 --version');
