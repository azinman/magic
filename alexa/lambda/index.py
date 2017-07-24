# Copyright Aaron Zinman 2017, 2018
# Copyright Duck Research LLC 2017, 2018
# All rights reserved.
#
# This file is part of Magichaus.
#
# Magichaus is free software: you can redistribute it and/or modify
# it under the terms of the GNU Lesser General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Magichaus is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public License
# along with Magichaus.  If not, see <http:#www.gnu.org/licenses/>.

"""
This sample demonstrates a simple skill built with the Amazon Alexa Skills Kit.
The Intent Schema, Custom Slots, and Sample Utterances for this skill, as well
as testing instructions are located at http://amzn.to/1LzFrj6

For additional samples, visit the Alexa Skills Kit Getting Started guide at
http://amzn.to/1LGWsLG
"""

from __future__ import print_function
import logging
import ssl
import uuid
import urllib
import urllib2

logger = logging.getLogger()
logger.setLevel(logging.DEBUG)

def lambda_handler(event, context):
    """ Route the incoming request based on type (LaunchRequest, IntentRequest,
    etc.) The JSON body of the request is provided in the event parameter.
    """
    logger.info('event: %s' % event)
    namespace = event['header']['namespace']
    if namespace == 'Alexa.ConnectedHome.Discovery':
        return handle_discovery(event)
    elif namespace == 'Alexa.ConnectedHome.Control':
        return handle_control(event)
    else:
        logger.error('No supported namespace: %s' % namespace)
        raise Exception('Something went wrong event: %s' % event)

def handle_discovery(event):
    # To later do something with
    accessToken = event['payload']['accessToken']

    headers = {
        'namespace': 'Alexa.ConnectedHome.Discovery',
        'name': 'DiscoverAppliancesResponse',
        'messageId': str(uuid.uuid4()),
        'payloadVersion': '2'
    }

    applianceDiscovered = {
        'applianceId': 'MagicHausCEC',
        'manufacturerName': 'MagicHaus',
        'modelName': 'M001',
        'version': 'VER01',
        'friendlyName': 'TV',
        'friendlyDescription': 'The TV in the Living Room',
        'isReachable': True,
        'actions': [
          'turnOn',
          'turnOff'
        ],
        'additionalApplianceDetails': {
            'fullApplianceId': '0db0155a-8d06-40aa-8faf-23c7a4ac1347'
        }
    }
    payloads = {
        'discoveredAppliances': [applianceDiscovered]
    }
    result = {
        'header': headers,
        'payload': payloads
    }
    logger.info('Discovery: %s' % result)
    return result

def handle_control(event):
    name = event['header']['name']
    accessToken = event['payload']['accessToken']
    appliance = event['payload']['appliance']
    details = appliance['additionalApplianceDetails']
    applianceUuid = details['fullApplianceId']
    applianceId = appliance['applianceId']

    if name == 'TurnOnRequest':
        logger.info('Turn on request for %s' % applianceId)
        responseName = 'TurnOnConfirmation'
        value = 'on'
    elif name == 'TurnOffRequest':
        logger.info('Turn off request for %s' % applianceId)
        responseName = 'TurnOffConfirmation'
        value = 'off'
    else:
        raise Exception('Unsupported operation')

    query = {'key': applianceUuid+'.power', 'value': value}
    encoded = urllib.urlencode(query)
    url = 'https://755.magicha.us:4000/put'
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    print(urllib2.urlopen(url, encoded, context=ctx).read())

    headers = {
        'namespace': 'Alexa.ConnectedHome.Control',
        'name': responseName,
        'messageId': str(uuid.uuid4()),
        'payloadVersion': '2'
    }
    payloads = {'success': True}
    result = {
        'header': headers,
        'payload': payloads
    }

    return result
