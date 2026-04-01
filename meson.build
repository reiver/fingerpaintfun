project('fingerpaintfun', 'c',
  version: '0.1.0',
  license: 'GPL-3.0-or-later',
)

go = find_program('go')

custom_target('fingerpaintfun',
  output: 'fingerpaintfun',
  command: [go, 'build', '-v', '-mod=vendor', '-o', '@OUTPUT@', '.'],
  build_by_default: true,
  install: true,
  install_dir: get_option('bindir'),
)

install_data(
  'data/fun.fingerpaint.app.desktop.in',
  rename: 'fun.fingerpaint.app.desktop',
  install_dir: join_paths(get_option('datadir'), 'applications'),
)
