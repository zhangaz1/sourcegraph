import { boolean } from '@storybook/addon-knobs'
import { storiesOf } from '@storybook/react'
import React from 'react'
import { LoaderInput } from './LoaderInput'
import { WebStory } from './WebStory'

const { add } = storiesOf('web/LoaderInput', module).addDecorator(story => (
    <div className="container mt-3" style={{ width: 800 }}>
        {story()}
    </div>
))

add('Interactive', () => (
    <WebStory>
        {() => (
            <>
                <LoaderInput loading={boolean('loading', true)}>
                    <input type="text" placeholder="Loader input" className="form-control" />
                </LoaderInput>
            </>
        )}
    </WebStory>
))
