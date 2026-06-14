const { execSync } = require('child_process');

const REPO_URL = 'https://rifaterdemsahin.github.io/claude-architect-certification/';
const RENDERER = 'markdown_renderer.html?file=';

function run(command) {
    console.log(`> Executing: ${command}`);
    try {
        return execSync(command, { encoding: 'utf8' });
    } catch (error) {
        console.error(`Error executing command: ${command}`);
        console.error(error.stderr || error.message);
        process.exit(1);
    }
}

function getLiveUrl(filePath) {
    if (filePath.endsWith('.md')) {
        return `${REPO_URL}${RENDERER}${filePath}`;
    }
    return `${REPO_URL}${filePath}`;
}

async function main() {
    const args = process.argv.slice(2);
    if (args.length === 0) {
        console.error('Usage: node batched_push.cjs \'<json_groups>\'');
        process.exit(1);
    }

    let groups;
    try {
        groups = JSON.parse(args[0]);
    } catch (e) {
        console.error('Failed to parse groups JSON.');
        process.exit(1);
    }

    for (let i = 0; i < groups.length; i++) {
        const { message, files } = groups[i];
        
        console.log(`\n📦 Batch ${i + 1}/${groups.length}: ${message}`);
        
        // Stage files
        files.forEach(file => {
            run(`git add "${file}"`);
        });

        // Commit
        run(`git commit -m "${message}"`);

        // Push
        run('git push');

        console.log('✅ Pushed successfully.');
        console.log('🔗 Active Links:');
        files.forEach(file => {
            console.log(`  - ${getLiveUrl(file)}`);
        });

        if (i < groups.length - 1) {
            console.log(`\n⏳ Waiting 5 minutes before next batch...`);
            // We use a synchronous sleep for simplicity in this script context
            // though in a real app we might use something else.
            // In Node, there isn't a built-in sync sleep besides busy-waiting or execSync('sleep')
            execSync('sleep 300');
        }
    }

    console.log('\n🚀 All batches completed!');
}

main();
