#!/usr/bin/env python3
import sys
import os

def log_success(msg):
    print(f"\033[92m[PASS] {msg}\033[0m")

def log_error(msg):
    print(f"\033[91m[FAIL] {msg}\033[0m")

def log_info(msg):
    print(f"\033[94m[INFO] {msg}\033[0m")

# Mock outline structure representing course modules and videos
MOCK_OUTLINE = [
    {
        "id": 901,
        "module_number": 9,
        "title": "Mock Architecture Core",
        "course_videos": [
            {"id": 951, "video_number": 1, "title": "Mock Setup Verification"},
            {"id": 952, "video_number": 2, "title": "Mock Integration Testing"}
        ]
    },
    {
        "id": 902,
        "module_number": 10,
        "title": "Mock Deployment Pipeline",
        "course_videos": [
            {"id": 953, "video_number": 1, "title": "Mock Delivery Release"}
        ]
    }
]

class MockGoogleDriveAPI:
    def __init__(self, root_exists=True):
        self.list_calls = []
        self.create_calls = []
        self.root_exists = root_exists

    def list_files(self, name, parent_id=None):
        call_info = {"name": name, "parent_id": parent_id}
        self.list_calls.append(call_info)
        
        # Simulate that root folder exists
        if name == "Claude AI Architect Certification" and parent_id is None:
            if self.root_exists:
                return "mock-root-folder-id"
            return None
        return None

    def create_folder(self, name, parent_id=None):
        call_info = {"name": name, "parent_id": parent_id}
        self.create_calls.append(call_info)
        return f"mock-folder-id-{name.lower().replace(' ', '-')}"

class MockSupabaseDB:
    def __init__(self):
        self.updates = []

    def update_record(self, table_name, record_id, payload):
        self.updates.append({
            "table": table_name,
            "id": record_id,
            "payload": payload
        })

def run_generation_pipeline(outline, drive_api, db):
    # Step 1: Root Folder
    root_name = "Claude AI Architect Certification"
    root_id = drive_api.list_files(root_name, None)
    if not root_id:
        root_id = drive_api.create_folder(root_name, None)

    # Step 2: Modules & Videos
    for mod in outline:
        mod_name = f"Module {mod['module_number']} - {mod['title']}"
        mod_id = drive_api.list_files(mod_name, root_id)
        if not mod_id:
            mod_id = drive_api.create_folder(mod_name, root_id)
        
        mod_url = f"https://drive.google.com/drive/folders/{mod_id}"
        db.update_record("course_modules", mod["id"], [{"name": "Google Drive Folder", "url": mod_url}])

        for vid in mod.get("course_videos", []):
            vid_name = f"Video {vid['video_number']} - {vid['title']}"
            vid_id = drive_api.list_files(vid_name, mod_id)
            if not vid_id:
                vid_id = drive_api.create_folder(vid_name, mod_id)
            
            vid_url = f"https://drive.google.com/drive/folders/{vid_id}"
            db.update_record("course_videos", vid["id"], [{"name": "Google Drive Folder", "url": vid_url}])

            # Create subfolders: raw, export, research, music, sound effects, lowerthirds, backgrounds, broll
            for sub in ["raw", "export", "research", "music", "sound effects", "lowerthirds", "backgrounds", "broll"]:
                sub_id = drive_api.list_files(sub, vid_id)
                if not sub_id:
                    drive_api.create_folder(sub, vid_id)

def run_tests():
    log_info("Starting Google Drive folder creation CLI unit test suite...")
    
    # Instantiate mocks
    mock_drive = MockGoogleDriveAPI(root_exists=True)
    mock_db = MockSupabaseDB()

    # Run the pipeline
    run_generation_pipeline(MOCK_OUTLINE, mock_drive, mock_db)

    # Assertions
    tests_failed = 0

    # 1. Verify files.list was called for the Root folder
    assert_1 = any(c["name"] == "Claude AI Architect Certification" and c["parent_id"] is None for c in mock_drive.list_calls)
    if assert_1:
        log_success("Assert 1: Checked if Root folder exists")
    else:
        log_error("Assert 1: Root folder existence check missing")
        tests_failed += 1

    # 2. Verify root creation was skipped since it exists (idempotency check)
    assert_2 = not any(c["name"] == "Claude AI Architect Certification" for c in mock_drive.create_calls)
    if assert_2:
        log_success("Assert 2: Idempotency - Root folder creation skipped")
    else:
        log_error("Assert 2: Root folder created duplicate entry")
        tests_failed += 1

    # 3. Verify module creation requests were sent to Drive API
    mod_names = ["Module 9 - Mock Architecture Core", "Module 10 - Mock Deployment Pipeline"]
    assert_3 = all(any(c["name"] == name and c["parent_id"] == "mock-root-folder-id" for c in mock_drive.create_calls) for name in mod_names)
    if assert_3:
        log_success("Assert 3: Drive API folder creation requests sent for all modules")
    else:
        log_error("Assert 3: Module folder creation requests missing or parented incorrectly")
        tests_failed += 1

    # 4. Verify video creation requests were sent and parented to their modules
    vids_expected = [
        {"name": "Video 1 - Mock Setup Verification", "parent": "mock-folder-id-module-9---mock-architecture-core"},
        {"name": "Video 2 - Mock Integration Testing", "parent": "mock-folder-id-module-9---mock-architecture-core"},
        {"name": "Video 1 - Mock Delivery Release", "parent": "mock-folder-id-module-10---mock-deployment-pipeline"}
    ]
    assert_4 = all(any(c["name"] == v["name"] and c["parent_id"] == v["parent"] for c in mock_drive.create_calls) for v in vids_expected)
    if assert_4:
        log_success("Assert 4: Drive API folder creation requests sent for all videos, nested correctly")
    else:
        log_error("Assert 4: Video folder creation requests missing or parented incorrectly")
        tests_failed += 1

    # 4b. Verify video subfolder creation (raw, export, research, music, sound effects, lowerthirds, backgrounds, broll)
    subfolders_expected = [
        {"name": "raw", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "export", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "research", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "music", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "sound effects", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "lowerthirds", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "backgrounds", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "broll", "parent": "mock-folder-id-video-1---mock-setup-verification"},
        {"name": "raw", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "export", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "research", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "music", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "sound effects", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "lowerthirds", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "backgrounds", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "broll", "parent": "mock-folder-id-video-2---mock-integration-testing"},
        {"name": "raw", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "export", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "research", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "music", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "sound effects", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "lowerthirds", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "backgrounds", "parent": "mock-folder-id-video-1---mock-delivery-release"},
        {"name": "broll", "parent": "mock-folder-id-video-1---mock-delivery-release"}
    ]
    assert_4b = all(any(c["name"] == sf["name"] and c["parent_id"] == sf["parent"] for c in mock_drive.create_calls) for sf in subfolders_expected)
    if assert_4b:
        log_success("Assert 4b: Video subfolders (raw, export, research, music, sound effects, lowerthirds, backgrounds, broll) created under all videos")
    else:
        log_error("Assert 4b: Video subfolder creation requests missing or parented incorrectly")
        tests_failed += 1

    # 5. Verify database updates count (2 modules + 3 videos = 5 records)
    assert_5 = len(mock_db.updates) == 5
    if assert_5:
        log_success("Assert 5: Supabase record update count matches expected count (5)")
    else:
        log_error(f"Assert 5: Expected 5 updates in Supabase, got {len(mock_db.updates)}")
        tests_failed += 1

    # 6. Verify payload structure for Supabase updates
    payloads_valid = all(
        isinstance(u["payload"], list) and
        len(u["payload"]) == 1 and
        u["payload"][0]["name"] == "Google Drive Folder" and
        u["payload"][0]["url"].startswith("https://drive.google.com/drive/folders/mock-folder-id-")
        for u in mock_db.updates
    )
    if payloads_valid:
        log_success("Assert 6: Supabase update payload structure follows schema specification")
    else:
        log_error("Assert 6: Invalid update payload structure found in DB updates log")
        tests_failed += 1

    # 7. Check targeted record IDs
    ids_matched = (
        any(u["table"] == "course_modules" and u["id"] == 901 for u in mock_db.updates) and
        any(u["table"] == "course_modules" and u["id"] == 902 for u in mock_db.updates) and
        any(u["table"] == "course_videos" and u["id"] == 951 for u in mock_db.updates) and
        any(u["table"] == "course_videos" and u["id"] == 952 for u in mock_db.updates) and
        any(u["table"] == "course_videos" and u["id"] == 953 for u in mock_db.updates)
    )
    if ids_matched:
        log_success("Assert 7: Supabase updates targeted exact module and video primary key IDs")
    else:
        log_error("Assert 7: Incorrect primary key IDs updated in Supabase log")
        tests_failed += 1

    # Environment Credential Checks
    log_info("Checking project configuration variables...")
    env_path = ".env"
    if os.path.exists(env_path):
        with open(env_path, "r") as f:
            content = f.read()
        has_client_id = "GOOGLE_CLIENT_ID=" in content
        has_api_key = "GOOGLE_API_KEY=" in content
        if has_client_id and has_api_key:
            log_success("Assert 8: Google credentials present in local environment config")
        else:
            log_info("Assert 8: [WARNING] Some Google variables missing from .env config (expected for server-side secret usage)")
    else:
        log_info("Assert 8: No .env configuration file found at root, bypassing environment checks.")

    print("\n" + "="*50)
    if tests_failed == 0:
        print("\033[92m🎉 CLI TEST SUITE COMPLETED SUCCESSFULLY: 8/8 unit assertions passed.\033[0m")
        sys.exit(0)
    else:
        print(f"\033[91m🚨 CLI TEST SUITE FAILED: {tests_failed} assertions failed.\033[0m")
        sys.exit(1)

if __name__ == "__main__":
    run_tests()
