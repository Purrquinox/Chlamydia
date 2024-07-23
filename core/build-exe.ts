import exe from '@angablue/exe';

const build = exe({
    entry: './coredist/api/api.js',
    out: './windows/chlamydia-core.exe',
    target: 'latest-win-x64'
});
