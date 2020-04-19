#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from .memory import Memory
from .pilot import Pilot
from . import utils

import sys

__version__ = '1.0.0'

print('using Pilot v{} ...'.format(__version__))

if sys.version_info.major < 3:
    msg = 'Pilot Requires Python 3.7 or greater. You are using {}'.format(
        sys.version)
    raise ValueError(msg)
