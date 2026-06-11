import os
import sys

# Add project root to sys.path so we can import src.utils.keyvault
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

try:
    from env_load import load_env  # custom helper if it exists, or just read .env
except ImportError:
    # simple dotenv loader fallback
    def load_env():
        env_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.env"))
        if os.path.exists(env_path):
            with open(env_path, "r") as f:
                for line in f:
                    if line.strip() and not line.startswith("#"):
                        key, val = line.strip().split("=", 1)
                        os.environ[key] = val.strip(' "')

load_env()

from 5_Symbols.src.utils.keyvault import get_secret

# Let's import get_secret. But wait, "5_Symbols" package name starts with a number.
# We can dynamically import it or use a simpler client.
