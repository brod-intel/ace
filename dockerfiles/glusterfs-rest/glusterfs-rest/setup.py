# -*- coding: utf-8 -*-
"""
    setup.py

    :copyright: (c) 2014 by Aravinda VK
    :license: MIT, see LICENSE for more details.
"""

from setuptools import setup


setup(
    name="GlusterFS REST API Server",
    version="0.2",
    packages=["glusterfsrest", "glusterfsrest.cli"],
    include_package_data=True,
    install_requires=['argparse==1.4.0', 'flask==1.1.2', 'gunicorn==20.0.4', 'pyyaml==5.4.1', 'jinja2==3.0.0a1', 'markupsafe==2.0.0a1', 'werkzeug==1.0.1', 'click==8.0.0a1', 'itsdangerous==2.0.0a1'],
    entry_points={
        "console_scripts": [
            "glusterrest = glusterfsrest.glusterrest:main",
        ]
    },
    package_data={'glusterfsrest': ['doc/*.yml', 'templates/*', 'static/*']},
    data_files=[('/usr/bin', ['bin/glusterrestd']),
                ('/var/lib/glusterd/rest/',
                 ['data/port'])],
    platforms="linux",
    zip_safe=False,
    author="Aravinda VK",
    author_email="mail@aravindavk.in",
    description="GlusterFS REST API server",
    license="BSD",
    keywords="glusterfs, cli, rest",
    url="https://github.com/aravindavk/glusterfs-rest",
)
